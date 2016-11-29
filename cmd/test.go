package main

import (
       "github.com/whosonfirst/go-whosonfirst-placetypes"
       "fmt"
)

func main() {

     s, _ := placetypes.Spec()
     fmt.Printf("%v\n", s)
}
