package net

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/PuerkitoBio/gred/cmd"
	"github.com/PuerkitoBio/gred/resp"
	"github.com/PuerkitoBio/gred/srv"
)

var (
	defdb = srv.NewDB(0)
)

// Conn defines the methods required to handle a connection to the
// server.
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
				return errors.New("db.Conn.Handle: write failed: " + err.Error())
			}
			continue
		}

		// Run the command
		var res interface{}
		var rerr error
		if cd, ok := cmd.Commands[strings.ToLower(ar[0])]; ok {
			args, ints, floats, err := cd.Parse(ar[0], ar[1:])
			if err != nil {
				rerr = err
			} else {
				switch cd := cd.(type) {
				case cmd.DBCmd:
					res, rerr = cd.ExecWithDB(c.db, args, ints, floats)
				case cmd.SrvCmd:
					res, rerr = cd.Exec(args, ints, floats)
				default:
					panic(fmt.Sprintf("unsupported command type: %T", cd))
				}
			}
		} else {
			rerr = fmt.Errorf("ERR unknown command '%s'", ar[0])
		}
		err = c.writeResponse(res, rerr)
		if err != nil {
			return err
		}
		if rerr == cmd.ErrQuit {
			return nil
		}
	}
}

// writeResponse writes the response to the network connection.
func (c *conn) writeResponse(res interface{}, err error) error {
	if err != nil {
		return resp.Encode(c, resp.Error(err.Error()))
	}
	return resp.Encode(c, res)
}
