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
	// ok is the standard OK serialized return value, to avoid allocation for
	// this common case.
	ok = []byte("+OK\r\n")

	// emptyArray is the unique empty array value used to return empty lists
	// for this common case.
	emptyArray = resp.Array{}

	// defaultDb is the default database used by all new connections (db 0).
	defaultDb = NewDB(0)

	// errInvalidCommand is returned when a malformed command is received.
	errInvalidCommand = errors.New("db: invalid command")

	// errNilSuccess is a sentinel value to indicate the success of a command,
	// and that the nil value should be returned.
	errNilSuccess = errors.New("db: (nil)")

	// errInvalidKeyType is the error returned when the key is not the right
	// type for the command to execute.
	errInvalidKeyType = errors.New("db: invalid key type")
)

// Conn represents a network connection to the server.
type Conn struct {
	c net.Conn

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
	case "expire":
		res, err = cmdExpire(c.ctx)
	case "persist":
		res, err = cmdPersist(c.ctx)
	case "ttl":
		res, err = cmdTTL(c.ctx)

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

		// Hashes commands
	case "hdel":
		res, err = cmdHdel(c.ctx)
	case "hexists":
		res, err = cmdHexists(c.ctx)
	case "hget":
		res, err = cmdHget(c.ctx)
	case "hgetall":
		res, err = cmdHgetAll(c.ctx)
	case "hkeys":
		res, err = cmdHkeys(c.ctx)
	case "hlen":
		res, err = cmdHlen(c.ctx)
	case "hmget":
		res, err = cmdHmget(c.ctx)
	case "hset":
		res, err = cmdHset(c.ctx)
	case "hvals":
		res, err = cmdHvals(c.ctx)

		// Lists commands
	case "lpop":
		res, err = cmdLpop(c.ctx)
	case "lpush":
		res, err = cmdLpush(c.ctx)
	case "rpop":
		res, err = cmdRpop(c.ctx)
	case "rpush":
		res, err = cmdRpush(c.ctx)

	default:
		err = errInvalidCommand
	}
	return c.writeResponse(res, err)
}

// writeResponse writes the response to the network connection.
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
