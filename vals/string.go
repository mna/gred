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
