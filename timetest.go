package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		t := time.Now()
		tmp := t.UnixNano() / int64(100000000) % int64(10)
		if tmp == 0 {
			fmt.Println(tmp)
		}
	}
}
