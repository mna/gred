package vals

type Value interface {
	Type() string
}

var empty = []string{}
