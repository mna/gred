package vals

type String interface {
	Value

	Append(string) int
	Get() string
	GetRange(int, int) string
	GetSet(string) string
	Set(string)
	StrLen() int
}

type stringval string

func (s *stringval) Type() string {
	return "string"
}

func (s *stringval) Append(v string) int {
	*s += stringval(v)
	return len(*s)
}

func (s *stringval) Get() string {
	return string(*s)
}

func (s *stringval) GetRange(start, end int) string {
	l := len(*s)
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

func (s *stringval) StrLen() int {
	return len(*s)
}

func NewString() String {
	return new(stringval)
}
