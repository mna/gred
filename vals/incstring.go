package vals

import "strconv"

// IncString defines the methods required to implement an incrementable
// string. IncString is a String with additional methods to increment
// or decrement a numeric value.
type IncString interface {
	String
	Decr() (int64, bool)
	DecrBy(int64) (int64, bool)
	Incr() (int64, bool)
	IncrBy(int64) (int64, bool)
	IncrByFloat(float64) (string, bool)
}

// Static type check to validate that *incString implements IncString.
var _ IncString = (*incString)(nil)

// incString implements an IncString.
type incString struct {
	String
	// TODO : May hold the parsed integer value eventually
}

// NewIncString creates a new IncString with the provided initial value.
func NewIncString(initval string) IncString {
	return &incString{
		NewString(initval),
	}
}

// Decr decrements the value by 1. It returns false as second return value
// if the current string value is not numeric.
func (is *incString) Decr() (int64, bool) {
	return is.IncrBy(-1)
}

// DecrBy decrements the value by dec. It returns false as second return value
// if the current string value is not numeric.
func (is *incString) DecrBy(dec int64) (int64, bool) {
	return is.IncrBy(-1 * dec)
}

// Incr increments the value by 1. It returns false as second return value
// if the current string value is not numeric.
func (is *incString) Incr() (int64, bool) {
	return is.IncrBy(1)
}

// IncrBy increments the value by inc. It returns false as second return value
// if the current string value is not numeric.
func (is *incString) IncrBy(inc int64) (int64, bool) {
	v, err := strconv.ParseInt(is.Get(), 10, 64)
	if err != nil {
		return 0, false
	}
	v += inc
	is.Set(strconv.FormatInt(v, 10))
	return v, true
}

// IncrByFloat increments the value by inc. It returns false as second return value
// if the current string value is not numeric.
func (is *incString) IncrByFloat(inc float64) (string, bool) {
	v, err := strconv.ParseFloat(is.Get(), 64)
	if err != nil {
		return "", false
	}
	v += inc
	// TODO: Limit to 17 digits precision, like Redis?
	ret := strconv.FormatFloat(v, 'f', -1, 64)
	is.Set(ret)
	return ret, true
}
