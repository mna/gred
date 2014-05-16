package srv

type Conn interface {
	Auth(string) bool
	Echo(string) string
	Ping() string
	Quit()
	Select(int) bool
}
