package types

// String defines the methods required to implement a String.
type String interface {
	Value

	Append(string) int64
	Get() string
	GetRange(int64, int64) string
	GetSet(string) string
	Set(string)
	SetRange(int64, string) int64
	StrLen() int64
}

// stringval is the internal representation of a String.
type stringval string

// NewString creates a new String holding the specified initial value.
func NewString(initval string) String {
	s := stringval(initval)
	return &s
}

// Type returns the type of the value, which is "string".
func (s *stringval) Type() string {
	return "string"
}

// Append appends the value v to the current string value.
// It returns the new length of the string.
func (s *stringval) Append(v string) int64 {
	*s += stringval(v)
	return int64(len(*s))
}

// Get returns the current string value.
func (s *stringval) Get() string {
	return string(*s)
}

// GetRange returns the value of the string from start to end.
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

// GetSet sets the value to v and returns the previous value.
func (s *stringval) GetSet(v string) string {
	old := *s
	*s = stringval(v)
	return string(old)
}

// Set sets the value to v.
func (s *stringval) Set(v string) {
	*s = stringval(v)
}

// SetRange sets a substring of the current value to v, starting
// at offset ofs. It returns the length of the new value.
func (s *stringval) SetRange(ofs int64, v string) int64 {
	// Fast path if there's no value to set
	if len(v) == 0 {
		return int64(len(*s))
	}

	// Pad with 0 bytes if required
	pad := int(ofs) + len(v) - len(*s)
	b := []byte(*s)
	if pad > 0 {
		b = append(b, make([]byte, pad)...)
	}

	// Set the new value in place
	copy(b[ofs:], v)
	*s = stringval(b)

	return int64(len(*s))
}

// StrLen returns the length of the string.
func (s *stringval) StrLen() int64 {
	return int64(len(*s))
}
