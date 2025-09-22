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

func ChildrenFunc() js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		pt_name := args[0].String()
		
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			pt, err := placetypes.GetPlacetypeByName(pt_name)
			
			if err != nil {
				reject.Invoke("Invalid placetype, %w", err)
				return nil
			}

			ch := placetypes.Children(pt)
			enc_ch, err := json.Marshal(ch)

			if err != nil {
				reject.Invoke("Failed to marshal response, %w", err)
				return nil
			}
			
			resolve.Invoke(string(enc_ch))
			return nil		
		})
			
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

func DescendantsFunc() js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		pt_name := args[0].String()
		roles := make([]string, 0)

		if len(args) > 1 {

			for _, a := range args[1:] {
				roles = append(roles, a.String())
			}
		}
		
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			pt, err := placetypes.GetPlacetypeByName(pt_name)
			
			if err != nil {
				reject.Invoke("Invalid placetype, %w", err)
				return nil
			}

			var d []*placetypes.WOFPlacetype
			
			if len(roles) > 1 {
				d = placetypes.DescendantsForRoles(pt, roles)
			} else {
				d = placetypes.Descendants(pt)
			}
			
			enc_d, err := json.Marshal(d)

			if err != nil {
				reject.Invoke("Failed to marshal response, %w", err)
				return nil
			}
			
			resolve.Invoke(string(enc_d))
			return nil		
		})
			
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

func AncestorsFunc() js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		pt_name := args[0].String()
		roles := make([]string, 0)

		if len(args) > 1 {

			for _, a := range args[1:] {
				roles = append(roles, a.String())
			}
		}
		
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			pt, err := placetypes.GetPlacetypeByName(pt_name)
			
			if err != nil {
				reject.Invoke("Invalid placetype, %w", err)
				return nil
			}

			var d []*placetypes.WOFPlacetype
			
			if len(roles) > 1 {
				d = placetypes.AncestorsForRoles(pt, roles)
			} else {
				d = placetypes.Ancestors(pt)
			}
			
			enc_d, err := json.Marshal(d)

			if err != nil {
				reject.Invoke("Failed to marshal response, %w", err)
				return nil
			}
			
			resolve.Invoke(string(enc_d))
			return nil		
		})
			
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

