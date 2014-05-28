package lists

import (
	"github.com/PuerkitoBio/gred/srv"
	"github.com/PuerkitoBio/gred/vals"
)

// unblock unblocks as many waiters as possible that are blocked waiting
// for a value from this key. Both the DB and the key must be under an
// exclusive lock.
func unblock(db srv.DB, k srv.Key, v vals.List) int {
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
