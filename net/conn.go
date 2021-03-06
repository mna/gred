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
	"github.com/golang/glog"
)

// NetConn defines the methods required to handle a network connection to the
// server.
type NetConn interface {
	io.ReadWriter
	Handle() error
}

var _ NetConn = (*netConn)(nil)
var _ srv.Conn = (*netConn)(nil)

// netConn represents a network connection to the server.
type netConn struct {
	net.Conn
	dbix int
}

// NewNetConn creates a new NetConn for the underlying net.Conn network
// connection.
func NewNetConn(c net.Conn) NetConn {
	conn := &netConn{
		Conn: c,
	}
	return conn
}

// Select sets the connection's DB index to ix.
func (c *netConn) Select(ix int) {
	c.dbix = ix
}

// Handle handles a connection to the server, and processes its requests.
func (c *netConn) Handle() error {
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

		if glog.V(2) {
			glog.Infof("[%s] command received: %v", c.RemoteAddr(), ar)
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
					// Get the connection's current database
					db, ok := srv.DefaultServer.GetDB(c.dbix)
					if !ok {
						panic(fmt.Sprintf("invalid database index: %d", c.dbix))
					}
					res, rerr = cd.ExecWithDB(db, args, ints, floats)
				case cmd.SrvCmd:
					res, rerr = cd.Exec(args, ints, floats)
				case cmd.ConnCmd:
					res, rerr = cd.ExecWithConn(c, args, ints, floats)
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
func (c *netConn) writeResponse(res interface{}, err error) error {
	if err != nil {
		if glog.V(2) {
			glog.Infof("[%s] response sent: %v", c.RemoteAddr(), err)
		}
		return resp.Encode(c, resp.Error(err.Error()))
	}
	if glog.V(2) {
		glog.Infof("[%s] response sent: %v", c.RemoteAddr(), res)
	}
	return resp.Encode(c, res)
}
