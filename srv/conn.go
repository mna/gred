package srv

import (
	"bufio"
	"errors"
	"io"
	"net"

	"github.com/PuerkitoBio/gred/resp"
)

type Conn interface {
	io.ReadWriter

	Handle() error
}

// conn represents a network connection to the server.
type conn struct {
	net.Conn

	db   DB
	quit bool
}

// NewConn creates a new Conn for the underlying net.Conn network
// connection.
func NewConn(c net.Conn) Conn {
	conn := &conn{
		Conn: c,
	}
	return conn
}

// Handle handles a connection to the server, and processes its requests.
func (c *conn) Handle() error {
	defer c.Close()

	br := bufio.NewReader(c)

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
		err = c.do(ar[0], ar[1:]...)
		if err != nil {
			return err
		}
		if c.quit {
			// Quit command, asked to close connection.
			return nil
		}
	}
}
