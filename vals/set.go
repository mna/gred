package vals

type Set interface {
	Value

	SAdd(...string) int64
	SCard() int64
	SDiff(...Set) []string
	SInter(...Set) []string
	SIsMember(string) bool
	SMembers() []string
	SPop() (string, bool)
	SRem(...string) int64
	SUnion(...Set) []string

	// Private method to get the internal map
	get() set
}

type set map[string]struct{}

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

func (s set) SCard() int64 {
	return int64(len(s))
}

func (s set) SDiff(vals ...Set) []string {
	var ret []string
	var ok bool

	for k := range s.get() {
		ok = true
		for _, other := range vals {
			otherSet := other.get()
			if _, ex := otherSet[k]; ex {
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

func (s set) SInter(vals ...Set) []string {
	var ret []string
	var ok bool

	for k := range s.get() {
		ok = true
		for _, other := range vals {
			otherSet := other.get()
			if _, ex := otherSet[k]; !ex {
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

func (s set) SIsMember(val string) bool {
	_, ok := s.get()[val]
	return ok
}

func (s set) SMembers() []string {
	ret := make([]string, len(s))
	i := 0
	for k := range s {
		ret[i] = k
		i++
	}
	return ret
}

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

func (s set) get() set {
	return s
}
