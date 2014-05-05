package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/gred/resp"
)

var (
	output = flag.String("o", "", "output file (empty outputs to stdout)")
	count  = flag.Int("n", 100, "number of SET commands to generate")
)

func main() {
	var err error

	flag.Parse()

	out := os.Stdout
	if *output != "" {
		out, err = os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer out.Close()

	ar := make(resp.Array, 3)
	for i := 0; i < *count; i++ {
		ar[0] = resp.BulkString("SET")
		ar[1] = resp.BulkString(fmt.Sprintf("key_%d", i))
		ar[2] = resp.BulkString(fmt.Sprintf("this is the string value for key #%d", i))
		err = resp.Encode(out, ar)
		if err != nil {
			log.Fatal(err)
		}
	}
}
