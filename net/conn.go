package net

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/gred/cmd"
	"github.com/PuerkitoBio/gred/resp"
	"github.com/PuerkitoBio/gred/srv"
)

var (
	pong  = []byte("+PONG\r\n")
	ok    = []byte("+OK\r\n")
	defdb = srv.NewDB(0)

	errArgNotInteger = errors.New("ERR value is not an integer or out of range")
	errArgNotFloat   = errors.New("ERR value is not a valid float")
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
			args, ints, floats, err := c.parseArgs(cd, ar[0], ar[1:])
			if err != nil {
				rerr = err
			} else {
				switch cd := cd.(type) {
				case cmd.DBCmd:
					res, rerr = cd.ExecWithDB(c.db, args, ints, floats)
				}
			}
		} else {
			rerr = fmt.Errorf("ERR unknown command '%s'", ar[0])
		}
		err = c.writeResponse(res, rerr)
		if err != nil {
			return err
		}
	}
}

func (c *conn) parseArgs(cd cmd.Cmd, name string, args []string) ([]string, []int, []float64, error) {
	l := len(args)
	ad := cd.GetArgDef()
	if l < ad.MinArgs || (l > ad.MaxArgs && ad.MaxArgs >= 0) {
		return nil, nil, nil, fmt.Errorf("ERR wrong number of arguments for '%s' command", name)
	}

	// Parse integers
	intix := ad.IntIndices
	ints := make([]int, len(intix))
	for i, ix := range intix {
		if ix < 0 {
			ix = l + ix
		}
		val, err := strconv.Atoi(args[ix])
		if err != nil {
			return nil, nil, nil, errArgNotInteger
		}
		ints[i] = val
	}

	// Parse floats
	fix := ad.FloatIndices
	floats := make([]float64, len(fix))
	for i, ix := range fix {
		if ix < 0 {
			ix = l + ix
		}
		val, err := strconv.ParseFloat(args[ix], 64)
		if err != nil {
			return nil, nil, nil, errArgNotFloat
		}
		floats[i] = val
	}
	return args, ints, floats, nil
}

// writeResponse writes the response to the network connection.
func (c *conn) writeResponse(res interface{}, err error) error {
	switch err {
	case cmd.ErrNilSuccess:
		// Special-case for success but nil return value
		return resp.Encode(c, nil)

	case cmd.ErrPong:
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
