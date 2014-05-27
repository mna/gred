package vals

/*
type WaitChan <-chan chan<- [2]string

type ListWaiter interface {
	Next() (WaitChan, bool)
}

var _ List = (*blockList)(nil)

type blockList struct {
	List
	tag string
}

func NewBList(tag string) List {
	return &blockList{
		List: NewList(),
		tag:  tag,
	}
}

func (b *blockList) LPush(vals ...string) int64 {
	// Perform the push
	ret := b.List.LPush(vals...)
	// Unblock any waiting subscriber
	b.unblock()
	return ret
}

func (b *blockList) RPush(vals ...string) int64 {
	// Perform the push
	ret := b.List.RPush(vals...)
	// Unblock any waiting subscriber
	b.unblock()
	return ret
}

func (b *blockList) unblock() {
	// While there are channels waiting for values, and there are values...
	for len(b.chans) > 0 && b.LLen() > 0 {
		ch := b.chans[0]
		r := b.isr[0]
		sendch, ok := <-ch

		// Was the waiting channel closed? If not, send it a value.
		if ok {
			// Has to return a value, because LLen is checked first
			var v string
			if r {
				v, _ = b.RPop()
			} else {
				v, _ = b.LPop()
			}
			sendch <- [2]string{b.tag, v}
		}

		// Continue with next waiting channel
		b.chans = b.chans[1:]
		b.isr = b.isr[1:]
	}
}
*/
