package main

import (
	"fmt"
	"time"
)

func main() {
	a := make(chan int)
	b := make(chan int)
	if b == nil {
		fmt.Println("3333333333")
	}
	go func() { b <- 10 }()
	go func() { a <- 10 }()
	for {
		select {
		case tmp := <-a:
			fmt.Println(tmp)
		case bb := <-b:
			fmt.Println("22222222", bb)
		case <-time.After(10 * time.Second):
			fmt.Println("11111111")
		}
	}
	<-time.After(3 * time.Second)
}
