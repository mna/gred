package vals

import "strconv"

// IncHash defines the methods required for an incrementable hash value.
// A IncHash is a Hash with additional methods to increment numeric values.
type IncHash interface {
	Hash
	HIncrBy(string, int64) (int64, bool)
	HIncrByFloat(string, float64) (string, bool)
}

// Static type check to validate that *incHash implements IncHash.
var _ IncHash = (*incHash)(nil)

// incHash implements a IncHash.
type incHash struct {
	Hash
	// TODO : May hold the parsed integer value eventually
}

// NewIncHash creates a new IncHash.
func NewIncHash() IncHash {
	return &incHash{
		NewHash(),
	}
}

// HIncrBy increments the value of field by inc. It creates the field
// and sets it at 0 before incrementing if field does not exist in the hash.
// It returns false as second return value if it could not perform the
// increment because the current value is not numeric.
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

// HIncrByFloat increments the value of field by inc. It creates the field
// and sets it at 0 before incrementing if field does not exist in the hash.
// It returns false as second return value if it could not perform the
// increment because the current value is not numeric.
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
