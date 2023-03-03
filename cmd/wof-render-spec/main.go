package main

import (
	"flag"
	"log"

	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-placetypes/draw"
)

func main() {

	path := flag.String("path", "placetypes.png", "...")

	flag.Parse()

	spec, err := placetypes.DefaultWOFPlacetypeSpecification()

	if err != nil {
		log.Fatalf("Failed to load specification, %v", err)
	}

	err = draw.DrawPlacetypesGraphToFile(spec, *path)

	if err != nil {
		log.Fatalf("Failed to draw placetypes graph, %v", err)
	}
}
