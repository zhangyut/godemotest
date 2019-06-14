package main

import (
	"fmt"
	"time"
)

func main() {
	var r chan []byte
	if r == nil {
		fmt.Println("is nil")
	} else {
		fmt.Println("not nil")
	}

	select {
	case r <- []byte("nini"):
		fmt.Println("nihoa")
	case <-time.After(10 * time.Second):
		fmt.Println("it's time")

	}
}
