package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/PuerkitoBio/gred/cmd"
	_ "github.com/PuerkitoBio/gred/cmd/connection"
	_ "github.com/PuerkitoBio/gred/cmd/hashes"
	_ "github.com/PuerkitoBio/gred/cmd/keys"
	_ "github.com/PuerkitoBio/gred/cmd/lists"
	_ "github.com/PuerkitoBio/gred/cmd/sets"
	_ "github.com/PuerkitoBio/gred/cmd/strings"
	gnet "github.com/PuerkitoBio/gred/net"
	"github.com/golang/glog"
)

// TODO : For optimization: http://confreaks.com/videos/3420-gophercon2014-building-high-performance-systems-in-go-what-s-new-and-best-practices

const (
	// port is the port to listen to
	port = 6379

	// maxSuccessiveConnErr is the maximum number of successive connection
	// errors before the server is stopped.
	maxSuccessiveConnErr = 3
)

func main() {
	flag.Parse()
	defer glog.Flush()

	// Print registered commands
	if glog.V(2) {
		for k := range cmd.Commands {
			glog.Infof("registered: %s", k)
		}
	}

	// Listen on TCP.
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	glog.V(1).Infof("listening on port %d", port)

	var errcnt int
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			errcnt++
			glog.Errorf("accept connection: %s", err)
			if errcnt >= maxSuccessiveConnErr {
				glog.Fatalf("%d successive connection errors, terminating...", errcnt)
			}
		}
		errcnt = 0
		glog.V(2).Infof("connection accepted: %s", conn.RemoteAddr())

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			conn := gnet.NewConn(c)
			err := conn.Handle()
			if err != nil {
				glog.Errorf("handle connection: %s", err)
			}
		}(conn)
	}
}
