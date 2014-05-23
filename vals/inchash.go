package vals

import "strconv"

type IncHash interface {
	Hash
	HIncrBy(string, int64) (int64, bool)
	HIncrByFloat(string, float64) (string, bool)
}

var _ IncHash = (*incHash)(nil)

type incHash struct {
	Hash
	// TODO : May hold the parsed integer value eventually
}

func NewIncHash() IncHash {
	return &incHash{
		NewHash(),
	}
}

func (ih *incHash) HIncrBy(field string, inc int64) (int64, bool) {
	val, ok := ih.HGet(field)
	if !ok {
		val = "0"
		ih.HSet(field, val)
	}
	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, false
	}
	v += inc
	ih.HSet(field, strconv.FormatInt(v, 10))
	return v, true
}

func (ih *incHash) HIncrByFloat(field string, inc float64) (string, bool) {
	val, ok := ih.HGet(field)
	if !ok {
		val = "0"
		ih.HSet(field, val)
	}
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return "", false
	}
	v += inc
	// TODO: Limit to 17 digits precision, like Redis?
	ret := strconv.FormatFloat(v, 'f', -1, 64)
	ih.HSet(field, ret)
	return ret, true
}
