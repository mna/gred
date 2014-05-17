package net

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/PuerkitoBio/gred/cmds"
	"github.com/PuerkitoBio/gred/resp"
	"github.com/PuerkitoBio/gred/srv"
)

var (
	pong  = []byte("+PONG\r\n")
	ok    = []byte("+OK\r\n")
	defdb = srv.NewDB(0)
)

type Conn interface {
	io.ReadWriter

	Handle() error
}

// conn represents a network connection to the server.
type conn struct {
	net.Conn

	db srv.DB
}

// NewConn creates a new Conn for the underlying net.Conn network
// connection.
func NewConn(c net.Conn) Conn {
	conn := &conn{
		Conn: c,
		db:   defdb,
	}
	return conn
}

// Handle handles a connection to the server, and processes its requests.
func (c *conn) Handle() error {
	defer c.Close()

	br := bufio.NewReader(c)

	var res interface{}
	for {
		// Get the request
		ar, err := resp.DecodeRequest(br)
		if err != nil {
			// Network error, return
			if _, ok := err.(net.Error); ok {
				return err
			}
			// Write the error to the client
			err = resp.Encode(c, resp.Error(err.Error()))
			if err != nil {
				// If write failed, return
				return errors.New("db.Conn.Handle: encode failed: " + err.Error())
			}
			continue
		}

		// Run the command
		if cmd, ok := cmds.Cmds[strings.ToLower(ar[0])]; ok {
			args, ints, floats := c.parseArgs(cmd, ar[1:])
			switch cmd := cmd.(type) {
			case cmds.DBCmd:
				res, err = cmd.ExecWithDB(c.db, args, ints, floats)
			}
			err = c.writeResponse(res, err)
		} else {
			err = c.writeResponse(nil, fmt.Errorf("ERR unknown command '%s'", ar[0]))
		}
		if err != nil {
			return err
		}
	}
}

func (c *conn) parseArgs(cmd cmds.Cmd, args []string) ([]string, []int, []float64) {
	return args, nil, nil
}

// writeResponse writes the response to the network connection.
func (c *conn) writeResponse(res interface{}, err error) error {
	switch err {
	case cmds.ErrNilSuccess:
		// Special-case for success but nil return value
		return resp.Encode(c, nil)

	case cmds.ErrPong:
		// Special-case for pong response
		_, err = c.Write(pong)
		return err

	case nil:
		if res == nil {
			// If the result is nil, send the OK response
			_, err = c.Write(ok)
			return err
		}
		// Otherwise encode the response
		return resp.Encode(c, res)

	default:
		// Return the non-nil error
		return resp.Encode(c, resp.Error(err.Error()))
	}
}
