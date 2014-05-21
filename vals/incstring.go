package vals

import "strconv"

type IncString interface {
	String
	Decr() (int64, bool)
	DecrBy(int64) (int64, bool)
	Incr() (int64, bool)
	IncrBy(int64) (int64, bool)
	IncrByFloat(float64) (string, bool)
}

var _ IncString = (*incString)(nil)

type incString struct {
	String
	// TODO : May hold the parsed integer value eventually
}

func (is *incString) Decr() (int64, bool) {
	return is.IncrBy(-1)
}

func (is *incString) DecrBy(dec int64) (int64, bool) {
	return is.IncrBy(-1 * dec)
}

func (is *incString) Incr() (int64, bool) {
	return is.IncrBy(1)
}

func (is *incString) IncrBy(inc int64) (int64, bool) {
	v, err := strconv.ParseInt(is.Get(), 10, 64)
	if err != nil {
		return 0, false
	}
	v += inc
	is.Set(strconv.FormatInt(v, 10))
	return v, true
}

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
