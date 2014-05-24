package vals

type List interface {
	Value

	LIndex(int64) (string, bool)
	LInsertBefore(string, string) int64
	LInsertAfter(string, string) int64
	LLen() int64
	LPop() (string, bool)
	LPush(...string) int64
	LRange(int64, int64) []string
	LRem(int64, string) int64
	LSet(int64, string) bool
	LTrim(int64, int64)
	RPop() (string, bool)
	RPush(...string) int64
}

const initListCap int = 10

var _ List = (*list)(nil)

type list []string

func NewList() List {
	l := make(list, 0, initListCap)
	return &l
}

func (l list) Type() string {
	return "list"
}

func (l *list) LIndex(ix int64) (string, bool) {
	ln := int64(len(*l))
	if ix < 0 {
		ix += ln
	}
	if ix >= 0 && ix < ln {
		return (*l)[ix], true
	}
	return "", false
}

func (l *list) LInsertBefore(pivot, val string) int64 {
	for i := 0; i < len(*l); i++ {
		if (*l)[i] == pivot {
			// Append a dummy value so that there's enough room in the slice
			*l = append(*l, "")
			// Copy all elements starting at pivot, one element to the right
			copy((*l)[i+1:], (*l)[i:])
			// Insert the new element at i
			(*l)[i] = val
			return int64(len(*l))
		}
	}
	return -1
}

func (l *list) LInsertAfter(pivot, val string) int64 {
	for i := 0; i < len(*l); i++ {
		if (*l)[i] == pivot {
			// Append a dummy value so that there's enough room in the slice
			*l = append(*l, "")
			// Copy all elements starting at pivot, one element to the right
			copy((*l)[i+2:], (*l)[i+1:])
			// Insert the new element at i+1
			(*l)[i+1] = val
			return int64(len(*l))
		}
	}
	return -1
}

func (l *list) LLen() int64 {
	return int64(len(*l))
}

func (l *list) LPop() (string, bool) {
	if len(*l) == 0 {
		return "", false
	}
	val := (*l)[0]
	(*l)[0] = ""
	*l = (*l)[1:]
	return val, true
}

func (l *list) LPush(vals ...string) int64 {
	// Reverse vals, then append
	for i, j := 0, len(vals)-1; i < j; i, j = i+1, j-1 {
		vals[i], vals[j] = vals[j], vals[i]
	}
	*l = append(vals, *l...)
	return int64(len(*l))
}

func (l *list) LRange(start, stop int64) []string {
	start, stop = l.normalizeStartStop(start, stop)
	if stop-start < 0 {
		return empty
	}
	return (*l)[start : stop+1]
}

func (l *list) LRem(cnt int64, val string) int64 {
	var n int64
	ln := len(*l)
	if cnt >= 0 {
		for i := 0; i < ln; i++ {
			if (*l)[i] == val {
				l.del(i)
				i--
				ln--
				n++
				if cnt > 0 && n >= cnt {
					break
				}
			}
		}
	} else {
		cnt *= -1
		for i := ln - 1; i >= 0; i-- {
			if (*l)[i] == val {
				l.del(i)
				i--
				n++
				if n >= cnt {
					break
				}
			}
		}
	}
	return n
}

func (l *list) del(ix int) {
	copy((*l)[ix:], (*l)[ix+1:])
	(*l)[len(*l)-1] = ""
	*l = (*l)[:len(*l)-1]
}

func (l *list) LSet(ix int64, val string) bool {
	ln := int64(len(*l))
	if ix < 0 {
		ix += ln
	}
	if ix >= 0 && ix < ln {
		(*l)[ix] = val
		return true
	}
	return false
}

// LTrim returns the number of elements remaining in the list, so that
// the containing key can know if it should be deleted.
func (l *list) LTrim(start, stop int64) {
	start, stop = l.normalizeStartStop(start, stop)
	if stop-start < 0 {
		*l = (*l)[:0]
		return
	}
	*l = (*l)[start : stop+1]
}

func (l *list) RPop() (string, bool) {
	if len(*l) == 0 {
		return "", false
	}
	val := (*l)[len(*l)-1]
	(*l)[len(*l)-1] = ""
	*l = (*l)[:len(*l)-1]
	return val, true
}

func (l *list) RPush(vals ...string) int64 {
	*l = append(*l, vals...)
	return int64(len(*l))
}

func (l *list) normalizeStartStop(start, stop int64) (int64, int64) {
	ln := int64(len(*l))
	if start < 0 {
		start += ln
	}
	if stop < 0 {
		stop += ln
	}
	if start < 0 {
		start = 0
	}
	if stop >= ln {
		stop = ln - 1
	}
	return start, stop
}
