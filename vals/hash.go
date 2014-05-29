package vals

// Hash defines the methods required to implement the Hash Redis type.
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

// Static check to make sure hash implements Hash.
var _ Hash = hash(nil)

// hash is the internal implementation of the Hash interface.
type hash map[string]string

// NewHash creates a new Hash value.
func NewHash() Hash {
	return make(hash)
}

// Type returns the type of the value, which is "hash".
func (h hash) Type() string {
	return "hash"
}

// HDel deletes the specified fields from the hash, and returns the
// number of fields removed.
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

// HExists returns true if the specified field exists in the hash.
func (h hash) HExists(field string) bool {
	_, ok := h[field]
	return ok
}

// HGet returns the value of the specified field. The second return value
// indicates if the field exists in the hash.
func (h hash) HGet(field string) (string, bool) {
	v, ok := h[field]
	return v, ok
}

// HGetAll returns the list of all key-value pairs in the hash.
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

// HKeys returns the list of keys in the hash.
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

// HLen returns the number of fields in the hash.
func (h hash) HLen() int64 {
	return int64(len(h))
}

// HMGet returns the list of values for all requested fields, in the
// order of the requested fields. It returns a nil value at the position
// of non-existing fields.
func (h hash) HMGet(fields ...string) []interface{} {
	ret := make([]interface{}, len(fields))
	for i, f := range fields {
		if v, ok := h[f]; ok {
			ret[i] = v
		}
	}
	return ret
}

// HMSet sets the values for all key-value tuples as received as argument.
func (h hash) HMSet(tuples ...string) {
	for i := 0; i < len(tuples); {
		h[tuples[i]] = tuples[i+1]
		i += 2
	}
}

// HSet sets the value of field to val, and returns true if the field had to be
// created.
func (h hash) HSet(field, val string) bool {
	_, ok := h[field]
	h[field] = val
	return !ok
}

// HSetNx sets the value of field to val only if the field does not already exists
// in the hash. It returns true if it did create and set the field.
func (h hash) HSetNx(field, val string) bool {
	if _, ok := h[field]; !ok {
		h[field] = val
		return true
	}
	return false
}

// HVals returns the list of values in the hash.
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
