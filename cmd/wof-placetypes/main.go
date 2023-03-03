package main

import (
	"log"

	"github.com/whosonfirst/go-whosonfirst-placetypes"
)

func main() {

	pt, err := placetypes.Placetypes()

	if err != nil {
		log.Fatal(err)
	}

	for i, p := range pt {
		log.Println(i, p.Name)
	}
}
