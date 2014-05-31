package vals

// Set defines the methods required to implement a Set.
type Set interface {
	Value

	SAdd(...string) int64
	SCard() int64
	SDiff(...Set) []string
	SInter(...Set) []string
	SIsMember(string) bool
	SMembers() []string
	SRem(...string) int64
	SUnion(...Set) []string
}

// Static type check to validate that *set implements Set.
var _ Set = (*set)(nil)

// set is the internal implementation of a Set.
type set map[string]struct{}

// NewSet creates a new Set.
func NewSet() Set {
	return make(set)
}

// Type returns the type of the value, which is "set".
func (s set) Type() string {
	return "set"
}

// SAdd adds the values to the set. It returns the number of values
// that were actually added.
func (s set) SAdd(vals ...string) int64 {
	var cnt int64

	for _, v := range vals {
		if _, ok := s[v]; !ok {
			s[v] = struct{}{}
			cnt++
		}
	}
	return cnt
}

// SCard returns the number of elements in the set.
func (s set) SCard() int64 {
	return int64(len(s))
}

// SDiff returns the elements found in the set that are not
// found in the other sets specified by vals.
func (s set) SDiff(vals ...Set) []string {
	var ok bool

	ret := []string{}
	for k := range s {
		ok = true
		for _, other := range vals {
			if ex := other.SIsMember(k); ex {
				ok = false
				break
			}
		}
		if ok {
			ret = append(ret, k)
		}
	}
	return ret
}

// SInter returns the intersection of the sets.
func (s set) SInter(vals ...Set) []string {
	var ok bool

	ret := []string{}
	for k := range s {
		ok = true
		for _, other := range vals {
			if ex := other.SIsMember(k); !ex {
				ok = false
				break
			}
		}
		if ok {
			ret = append(ret, k)
		}
	}
	return ret
}

// SIsMember returns true if the value val is in the set.
func (s set) SIsMember(val string) bool {
	_, ok := s[val]
	return ok
}

// SMembers returns the list of all members of the set.
func (s set) SMembers() []string {
	ret := make([]string, len(s))
	i := 0
	for k := range s {
		ret[i] = k
		i++
	}
	return ret
}

// SRem removes the values vals from the set. It returns the number
// of elements that were actually removed.
func (s set) SRem(vals ...string) int64 {
	var cnt int64
	for _, v := range vals {
		if _, ok := s[v]; ok {
			delete(s, v)
			cnt++
		}
	}
	return cnt
}

// SUnion returns the union of all sets.
func (s set) SUnion(sets ...Set) []string {
	ret := make(set, len(s))
	for k := range s {
		ret[k] = struct{}{}
	}
	for _, otherSet := range sets {
		m := otherSet.SMembers()
		for _, k := range m {
			ret[k] = struct{}{}
		}
	}
	return ret.SMembers()
}
