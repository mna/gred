package cmds

type Cmd interface {
	IntArgIndices() []int
	FloatArgIndices() []int
	NumArgs() (int, int)
}
