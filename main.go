package main

import (
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
			buf := make([]byte, 1024)
			n, err := c.Read(buf)
			if err != nil {
				log.Println(err)
			}
			log.Println(string(buf), n)
			ar, err := resp.DecodeRequest(buf)
			if err != nil {
				log.Println(err)
			}
			log.Printf("%#v\n", ar)
			c.Write([]byte("+OK\r\n"))
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}
