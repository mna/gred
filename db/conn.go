package db

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"strings"

	"github.com/PuerkitoBio/gred/resp"
)

var (
	pong = []byte("+PONG\r\n")
	ok   = []byte("+OK\r\n")

	defaultDb = NewDB(0)

	errPong = errors.New("pong")
)

type connCmdFunc func(c *Conn, args ...string) (interface{}, error)

type connCmdDef struct {
	fn    connCmdFunc
	nArgs int
}

var connCmds = map[string]connCmdDef{
	"ping": connCmdDef{(*Conn).ping, 0},
	"echo": connCmdDef{(*Conn).echo, 1},
}

// Conn represents a network connection to the server.
type Conn struct {
	c  net.Conn
	db *Database
}

// NewConn creates a new Conn for the underlying net.Conn network
// connection.
func NewConn(c net.Conn) *Conn {
	return &Conn{
		c:  c,
		db: defaultDb,
	}
}

// Handle handles a connection to the server, and processes its requests.
func (c *Conn) Handle() error {
	defer c.c.Close()

	br := bufio.NewReader(c.c)

	for {
		// Get the request
		ar, err := resp.DecodeRequest(br)
		if err != nil {
			// Network error, return
			if err == io.EOF || err == io.ErrClosedPipe {
				return err
			}
			// Write the error to the client
			err = resp.Encode(c.c, resp.Error(err.Error()))
			if err != nil {
				// If write failed, return
				return errors.New("db.Conn.Handle: encode failed: " + err.Error())
			}
			continue
		}
		log.Printf("%#v\n", ar)

		// Run the command
		err = c.Do(ar[0], ar[1:]...)
		if err != nil {
			return err
		}
	}
}

func (c *Conn) ping(args ...string) (interface{}, error) {
	// Special case for ping, avoid allocation and return the pong predefined response.
	return nil, errPong
}

func (c *Conn) echo(args ...string) (interface{}, error) {
	return args[0], nil
}

// Do executes a given command on the connection.
func (c *Conn) Do(cmd string, args ...string) error {
	var res interface{}
	var err error

	cmd = strings.ToLower(cmd)
	if def, ok := connCmds[cmd]; ok {
		if len(args) != def.nArgs {
			err = ErrMissingArg
		} else {
			res, err = def.fn(c, args...)
		}
	} else {
		res, err = c.db.Do(cmd, args...)
	}

	switch err {
	case errNilSuccess:
		// Special-case for success but nil return value
		return resp.Encode(c.c, nil)

	case errPong:
		// Special-case for pong response
		_, err = c.c.Write(pong)
		return err

	case nil:
		if res == nil {
			// If the result is nil, send the OK response
			_, err = c.c.Write(ok)
			return err
		}
		// Otherwise encode the response
		return resp.Encode(c.c, res)

	default:
		// Return the non-nil error
		return resp.Encode(c.c, resp.Error(err.Error()))
	}
}
