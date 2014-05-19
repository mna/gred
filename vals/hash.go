package vals

var empty = []string{}

type Hash interface {
	Value

	HDel(...string) int64
	HExists(string) bool
	HGet(string) (string, bool)
	HGetAll() []string
	HKeys() []string
	HLen() int64
	HMGet(...string) []interface{}
	HMSet(...string)
	HSet(string, string) bool
	HSetNx(string, string) bool
	HVals() []string
}

var _ Hash = hash(nil)

// hash is the internal implementation of the Hash interface.
type hash map[string]string

func NewHash() Hash {
	return make(hash)
}

func (h hash) Type() string {
	return "hash"
}

func (h hash) HDel(fields ...string) int64 {
	var cnt int64
	for _, f := range fields {
		if _, ok := h[f]; ok {
			cnt++
			delete(h, f)
		}
	}
	return cnt
}

func (h hash) HExists(field string) bool {
	_, ok := h[field]
	return ok
}

func (h hash) HGet(field string) (string, bool) {
	v, ok := h[field]
	return v, ok
}

func (h hash) HGetAll() []string {
	if len(h) == 0 {
		return empty
	}
	vals := make([]string, 2*len(h))
	i := 0
	for k, v := range h {
		vals[i] = k
		vals[i+1] = v
		i += 2
	}
	return vals
}

func (h hash) HKeys() []string {
	if len(h) == 0 {
		return empty
	}
	keys := make([]string, len(h))
	i := 0
	for k := range h {
		keys[i] = k
		i++
	}
	return keys
}

func (h hash) HLen() int64 {
	return int64(len(h))
}

func (h hash) HMGet(fields ...string) []interface{} {
	ret := make([]interface{}, len(fields))
	for i, f := range fields {
		if v, ok := h[f]; ok {
			ret[i] = v
		}
	}
	return ret
}

func (h hash) HMSet(tuples ...string) {
	for i := 0; i < len(tuples); {
		h[tuples[i]] = tuples[i+1]
		i += 2
	}
}

func (h hash) HSet(field, val string) bool {
	_, ok := h[field]
	h[field] = val
	return !ok
}

func (h hash) HSetNx(field, val string) bool {
	if _, ok := h[field]; !ok {
		h[field] = val
		return true
	}
	return false
}

func (h hash) HVals() []string {
	if len(h) == 0 {
		return empty
	}
	vals := make([]string, len(h))
	i := 0
	for _, v := range h {
		vals[i] = v
		i++
	}
	return vals
}
