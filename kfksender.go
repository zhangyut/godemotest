package main

import (
	"fmt"
	"quark/broadcast"
)

func main() {
	sender, err := broadcast.NewKfkSender("127.0.0.1:9092", "test", "127.0.0.1")
	if err != nil {
		panic(err)
	}

	err = sender.Send([]byte("我问你过的好不好"))
	if err != nil {
		panic(err)
	} else {
		fmt.Println("send succeed.")
	}
}
