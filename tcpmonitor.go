package main

import (
	"crystalball/prober/tcp"
	"flag"
	"fmt"
	"time"
)

var (
	ip      string
	port    string
	timeout int
)

func init() {
	flag.StringVar(&ip, "ip", "", "ip")
	flag.StringVar(&port, "port", "", "port")
	flag.IntVar(&timeout, "timeout", 3000, "timeout")
}

func main() {
	flag.Parse()
	if ip == "" || port == "" {
		panic("ip or port is nil.")
	}
	task := tcp.TcpConnTestTask{
		Id:      "tcp_conn_test",
		Ip:      ip,
		Port:    port,
		Timeout: timeout,
	}

	task.Validate()
	connTest := tcp.NewTcpConnTest()
	result := connTest.DoProbe(&task)
	fmt.Printf("%v, %s:%s result:%d\n", time.Now(), ip, port, result.GetRtt())
}
