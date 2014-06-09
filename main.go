package main

import (
	"flag"
	"log"
	"net"

	"github.com/PuerkitoBio/gred/cmd"
	_ "github.com/PuerkitoBio/gred/cmd/connection"
	_ "github.com/PuerkitoBio/gred/cmd/hashes"
	_ "github.com/PuerkitoBio/gred/cmd/keys"
	_ "github.com/PuerkitoBio/gred/cmd/lists"
	_ "github.com/PuerkitoBio/gred/cmd/server"
	_ "github.com/PuerkitoBio/gred/cmd/sets"
	_ "github.com/PuerkitoBio/gred/cmd/strings"
	gnet "github.com/PuerkitoBio/gred/net"
	"github.com/golang/glog"
)

// TODO : For optimization ideas: http://confreaks.com/videos/3420-gophercon2014-building-high-performance-systems-in-go-what-s-new-and-best-practices

const (
	// maxSuccessiveConnErr is the maximum number of successive connection
	// errors before the server is stopped.
	maxSuccessiveConnErr = 3
)

var (
	addr  = flag.String("addr", ":6379", "network address to listen to")
	iface = flag.String("net", "tcp", "network interface to use")
)

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()

	flag.Parse()
	defer glog.Flush()

	// Print registered commands
	if glog.V(2) {
		for k := range cmd.Commands {
			glog.Infof("registered: %s", k)
		}
	}

	l, err := net.Listen(*iface, *addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	glog.V(1).Infof("listening on %s://%s", *iface, *addr)

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
			conn := gnet.NewNetConn(c)
			err := conn.Handle()
			if err != nil {
				glog.Errorf("handle connection: %s", err)
			}
		}(conn)
	}
}
