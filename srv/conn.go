package srv

// Conn defines the methods required to implement a Connection.
type Conn interface {
	Select(int)
}
