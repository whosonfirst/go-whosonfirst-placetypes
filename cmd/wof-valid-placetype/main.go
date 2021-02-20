package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"log"
)

func main() {

	flag.Parse()

	for _, pt := range flag.Args() {
		log.Printf("%s\t%t\n", pt, placetypes.IsValidPlacetype(pt))
	}
}
