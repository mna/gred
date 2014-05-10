package db

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strings"

	"github.com/PuerkitoBio/gred/resp"
)

var (
	ok        = []byte("+OK\r\n")
	defaultDb = NewDB(0)

	// errInvalidCommand is returned when a malformed command is received.
	errInvalidCommand = errors.New("db: invalid command")

	// ErrNilSuccess is a sentinel value to indicate the success of a command,
	// and that the nil value should be returned.
	errNilSuccess = errors.New("db: (nil)")
)

// Conn represents a network connection to the server.
type Conn struct {
	c    net.Conn
	db   *Database
	ctx  *Ctx
	quit bool
}

// NewConn creates a new Conn for the underlying net.Conn network
// connection.
func NewConn(c net.Conn) *Conn {
	conn := &Conn{
		c:   c,
		db:  defaultDb,
		ctx: &Ctx{},
	}
	conn.ctx.conn = conn
	return conn
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

// Do executes a given command on the connection.
func (c *Conn) do(cmd string, args ...string) error {
	var res interface{}
	var err error

	// Prepare the command and context
	cmd = strings.ToLower(cmd)
	c.ctx.db = c.db
	c.ctx.raw = args

	// Excecute the command
	switch cmd {
	// Connection commands
	case "echo":
		res, err = cmdEcho(c.ctx)
	case "ping":
		res, err = cmdPing(c.ctx)
	case "quit":
		res, err = cmdQuit(c.ctx)
		if err == nil {
			c.quit = true
		}

		// Keys commands
	case "del":
		res, err = cmdDel(c.ctx)
	case "exists":
		res, err = cmdExists(c.ctx)

		// Strings commands
	case "append":
		res, err = cmdAppend(c.ctx)
	case "get":
		res, err = cmdGet(c.ctx)
	case "getrange":
		res, err = cmdGetRange(c.ctx)
	case "getset":
		res, err = cmdGetSet(c.ctx)
	case "set":
		res, err = cmdSet(c.ctx)
	case "strlen":
		res, err = cmdStrLen(c.ctx)

	default:
		err = errInvalidCommand
	}
	return c.writeResponse(res, err)
}

func (c *Conn) writeResponse(res interface{}, err error) error {
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
