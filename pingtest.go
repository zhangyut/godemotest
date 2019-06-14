package main

import (
	"crystalball/prober/ping"
	"fmt"
	"time"
)

func main() {
	pinger, _ := ping.NewSyncPinger()
	timeout := pinger.Ping("202.173.9.11", 5000)
	stimeout := fmt.Sprintf("timeout : %d", timeout)
	fmt.Println(stimeout)

	timeout = pinger.Ping("www.baidu.com", 5000)
	stimeout = fmt.Sprintf("timeout : %d", timeout)
	fmt.Println(stimeout)
	<-time.After(2 * time.Second)
}
