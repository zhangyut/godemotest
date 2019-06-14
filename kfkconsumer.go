package main

import (
	"fmt"
	"os"
	"os/signal"
	"quark/broadcast"
	"time"
)

func main() {
	receiver, err := broadcast.NewKfkReceiver("202.173.9.45:9092,202.173.9.47:9092,202.173.9.48:9092", "gslb2zcloud", "192.168.219.163")
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	for {
		select {
		case msg := <-receiver.GetDataChan():
			fmt.Println(string(msg))
		case <-time.After(1 * time.Minute):
			fmt.Println("timeout 1 minute.")
		case <-signals:
			fmt.Println("go out now.")
			goto goout
		}
	}
goout:
}
