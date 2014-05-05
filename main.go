package main

import (
	"log"
	"net"

	"github.com/PuerkitoBio/gred/db"
)

func main() {
	// Listen on TCP port 6379 on all interfaces.
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			conn := db.NewConn(c)
			err := conn.Handle()
			if err != nil {
				log.Println("ERROR: ", err)
			}
		}(conn)
	}
}