package main

import (
	"cement/cron"
	"fmt"
	"sync"
	"time"
)

func wait(wg *sync.WaitGroup) chan bool {
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	c := cron.New()
	c.AddFunc("* * * * * ?", func() { fmt.Println("run func.") })
	c.AddFunc("1 * * * * ?", func() { wg.Done() })
	c.Start()
	defer c.Stop()

	for {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("***")
		case <-wait(wg):
			return
		}
	}
}
