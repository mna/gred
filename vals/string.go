package vals

type String interface {
	Value

	Append(string) int64
	Get() string
	GetRange(int64, int64) string
	GetSet(string) string
	Set(string)
	StrLen() int64
}

type stringval string

func NewString(initval string) String {
	s := stringval(initval)
	return &s
}

func (s *stringval) Type() string {
	return "string"
}

func (s *stringval) Append(v string) int64 {
	*s += stringval(v)
	return int64(len(*s))
}

func (s *stringval) Get() string {
	return string(*s)
}

func (s *stringval) GetRange(start, end int64) string {
	l := int64(len(*s))
	if start < 0 {
		start = l + start
		if start < 0 {
			start = 0
		}
	}
	if start >= l {
		return ""
	}
	if end < 0 {
		end = l + end
	}
	if end < 0 || end < start {
		return ""
	}
	if end >= l {
		end = l - 1
	}
	return string((*s)[start : end+1])
}

func (s *stringval) GetSet(v string) string {
	old := *s
	*s = stringval(v)
	return string(old)
}

func (s *stringval) Set(v string) {
	*s = stringval(v)
}

func (s *stringval) StrLen() int64 {
	return int64(len(*s))
}
