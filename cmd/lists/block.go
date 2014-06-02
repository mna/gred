package lists

import (
	"time"

	"github.com/PuerkitoBio/gred/cmd"
	"github.com/PuerkitoBio/gred/srv"
	"github.com/PuerkitoBio/gred/types"
)

// unblock unblocks as many waiters as possible that are blocked waiting
// for a value from this key. Both the DB and the key must be under an
// exclusive lock.
func unblock(db srv.DB, k srv.Key, v types.List) int {
	var cnt int

	// While there are values...
	for v.LLen() > 0 {
		// Get a waiter
		ch, r := db.NextWaiter(k.Name())
		if ch == nil {
			// No more waiter, return
			return cnt
		}
		sendch, ok := <-ch

		// Was the waiting channel closed? If not, send it a value.
		if ok {
			// Has to return a value, because LLen is checked first
			var val string
			if r {
				val, _ = v.RPop()
			} else {
				val, _ = v.LPop()
			}
			cnt++
			sendch <- [2]string{k.Name(), val}
		}
	}
	return cnt
}

func blockPop(db srv.DB, secs int64, rpop bool, lists ...string) ([]string, error) {
	db.Lock()
	unlocks := make([]func(), 0)
	unlocks = append(unlocks, db.Unlock)

	keys := db.Keys()
	for _, nm := range lists {
		k, ok := keys[nm]
		// Ignore non-existing keys in non-blocking portion
		if ok {
			// Lock the key
			k.Lock()
			unlocks = append(unlocks, k.Unlock)

			// Get the value, if possible
			v := k.Val()
			if v, ok := v.(types.List); ok {
				var val string
				if rpop {
					val, ok = v.RPop()
				} else {
					val, ok = v.LPop()
				}
				if ok {
					// Delete the key if there are no more values
					if v.LLen() == 0 {
						db.DelKey(k.Name())
					}

					// Unlock all keys in reverse order, and return
					for i := len(unlocks) - 1; i >= 0; i-- {
						unlocks[i]()
					}
					return []string{k.Name(), val}, nil
				}
			} else {
				// Unlock all keys in reverse order
				for i := len(unlocks) - 1; i >= 0; i-- {
					unlocks[i]()
				}
				// Return invalid type error
				return nil, cmd.ErrInvalidValType
			}
		}
	}

	// If no value was readily available, now all keys are locked, enter
	// the waiting workflow.
	ch := make(chan chan<- [2]string)
	for _, nm := range lists {
		if rpop {
			db.WaitRPop(nm, ch)
		} else {
			db.WaitLPop(nm, ch)
		}
	}

	// Prepare channels (timeout and receive values)
	var timeoutCh <-chan time.Time
	if secs > 0 {
		timeoutCh = time.After(time.Duration(secs) * time.Second)
	}
	recCh := make(chan [2]string)

	// Unlock all locks so that other connections can proceed
	for i := len(unlocks) - 1; i >= 0; i-- {
		unlocks[i]()
	}

	// Wait for a value
	select {
	case ch <- (chan<- [2]string)(recCh):
		close(ch)
		vals := <-recCh
		return vals[:], nil
	case <-timeoutCh:
		close(ch)
		return nil, nil
	}
}
