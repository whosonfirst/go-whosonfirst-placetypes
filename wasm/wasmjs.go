//go:build wasmjs
package wasm

import (
	"syscall/js"
	"encoding/json"
	"fmt"
	
	"github.com/whosonfirst/go-whosonfirst-placetypes"
)

func PlacetypesFunc() js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]
			
			pt, err := placetypes.Placetypes()
			
			if err != nil {
				reject.Invoke(fmt.Printf("Failed to derive placetypes, %v\n", err))
				return nil
			}
			
			enc_pt, err := json.Marshal(pt)
			
			if err != nil {
				reject.Invoke(fmt.Printf("Failed to encode placetypes, %v\n", err))
				return nil
			}
			
			resolve.Invoke(string(enc_pt))
			return nil		
		})
			
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

func IsValidPlacetypeFunc() js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		pt_name := args[0].String()
		
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			if !placetypes.IsValidPlacetype(pt_name){
				reject.Invoke("Invalid placetype")
				return nil
			}
						
			resolve.Invoke()
			return nil		
		})
			
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

// Children

// Descendants

// Ancestors
