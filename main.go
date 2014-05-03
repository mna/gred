package main

import (
	"io"
	"log"
	"net"

	"github.com/PuerkitoBio/gred/resp"
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
			log.Println("request")
			for {
				ar, err := resp.DecodeRequest(c)
				if err != nil {
					log.Println(err)
					if err == io.EOF {
						break
					}
				}
				log.Printf("%s\n", ar)
				_, err = c.Write([]byte("+OK\r\n"))
				if err != nil {
					log.Println(err)
					if err == io.ErrClosedPipe {
						break
					}
				}
			}
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}
