package db

import (
	"errors"
	"io"
	"net"
	"strings"

	"github.com/PuerkitoBio/gred/resp"
)

var (
	pong = []byte("+PONG\r\n")
	ok   = []byte("+OK\r\n")

	defaultDb = NewDB(0)
)

type Conn struct {
	c  net.Conn
	db *Database
}

func NewConn(c net.Conn) *Conn {
	return &Conn{
		c:  c,
		db: defaultDb,
	}
}

func (c *Conn) Handle() error {
	defer c.c.Close()

	for {
		// Get the request
		ar, err := resp.DecodeRequest(c.c)
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

		// Run the command
		err = c.Do(ar[0], ar[1:]...)
		if err != nil {
			return err
		}
	}
}

func (c *Conn) Do(cmd string, args ...string) error {
	var res interface{}
	var err error

	switch lc := strings.ToLower(cmd); lc {
	case "ping":
		// Special case for ping, avoid allocation and return the pong predefined response.
		_, err = c.c.Write(pong)
		return err

	default:
		res, err = c.db.Do(lc, args...)
	}

	// If the command returned an error, send it back to the client
	if err != nil {
		// Special-case for success but nil return value
		if err == ErrNilSuccess {
			return resp.Encode(c.c, nil)
		}
		return resp.Encode(c.c, resp.Error(err.Error()))
	}
	// If the result is nil, send the OK response
	if res == nil {
		_, err = c.c.Write(ok)
		return err
	}
	// Otherwise encode the response
	return resp.Encode(c.c, res)
}
