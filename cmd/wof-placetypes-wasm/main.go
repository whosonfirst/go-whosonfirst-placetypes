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

	children_func := wasm.ChildrenFunc()
	defer children_func.Release()

	descendants_func := wasm.DescendantsFunc()
	defer descendants_func.Release()

	ancestors_func := wasm.AncestorsFunc()
	defer ancestors_func.Release()
	
	js.Global().Set("wof_placetypes", placetypes_func)
	js.Global().Set("wof_placetypes_is_valid", isvalid_func)
	js.Global().Set("wof_placetypes_children", children_func)
	js.Global().Set("wof_placetypes_descendants", descendants_func)
	js.Global().Set("wof_placetypes_ancestors", ancestors_func)				

	c := make(chan struct{}, 0)

	log.Println("wof_placetypes functions initialized")
	<-c
}

