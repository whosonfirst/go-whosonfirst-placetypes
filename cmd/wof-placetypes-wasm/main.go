//go:build wasmjs
package main

import (
	"log"
	"syscall/js"

	"github.com/whosonfirst/go-whosonfirst-placetypes/wasm"	
)

func main() {

	placetypes_func := wasm.PlacetypesFunc()
	defer placetypes_func.Release()

	isvalid_func := wasm.IsValidPlacetypeFunc()
	defer isvalid_func.Release()
	
	js.Global().Set("wof_placetypes", placetypes_func)
	js.Global().Set("wof_is_valid_placetype", isvalid_func)	

	c := make(chan struct{}, 0)

	log.Println("wof_placetypes functions initialized")
	<-c
}

