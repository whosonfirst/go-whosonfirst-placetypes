package main

import (
	"flag"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"log"
)

func main() {

	var roles multi.MultiString
	flag.Var(&roles, "role", "...")

	flag.Parse()

	for _, str_pt := range flag.Args() {

		pt, err := placetypes.GetPlacetypeByName(str_pt)

		if err != nil {
			log.Fatal(err)
		}

		var descendants []*placetypes.WOFPlacetype

		if len(roles) > 0 {
			descendants = placetypes.DescendantsForRoles(pt, roles)
		} else {
			descendants = placetypes.Descendants(pt)
		}

		for i, p := range descendants {
			log.Println(i, p.Name)
		}
	}
}
