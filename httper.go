package main

import (
	"crystalball/prober/http"
	"flag"
	"fmt"
	"time"
)

var (
	t    int
	host string
	url  string
)

func init() {
	flag.IntVar(&t, "t", 1, "thread nums")
	flag.StringVar(&host, "h", "", "dst host")
	flag.StringVar(&url, "u", "", "dst url")
}

type Probe struct {
	Type   string
	Notify chan int
}

func NewProbe(t string, notify chan int) *Probe {
	var probe Probe
	probe.Type = t
	probe.Notify = notify
	return &probe
}

func (p *Probe) Run(no int) {
	index := 0
	begin := time.Now()
	httpTask := http.HttpTask{
		Id:         "aaaaaaa",
		Host:       host,
		Url:        url,
		Timeout:    5000,
		Method:     "GET",
		CreateTime: 0,
	}

	for {
		go func() {
			httper := http.NewFasterCurl()
			tmp := httper.DoProbe(&httpTask)
			result, _ := tmp.(*http.HttpResult)

			if result.Rtt >= 5000 {
				errMsg := fmt.Sprintf("seq:%d-%d, host:%s, url:%s, ttl:%d\n", no, index, host, url, result.Rtt)
				fmt.Printf("failed:%s\n", errMsg)
			} else {
				p.Notify <- 0
				index = index + 1
			}
			if index >= 100 {
				index = 0
				end := time.Now()
				dur := end.Sub(begin)
				msg := fmt.Sprintf("thread no:%d accomplish 100个http, dur:%d秒", no, dur/(1*time.Second))
				fmt.Printf("log :%s\n", msg)
				begin = time.Now()
			}
		}()
		<-time.After(1 * time.Second)
	}
}

func main() {
	flag.Parse()
	if host == "" || url == "" {
		fmt.Println("host or url is nil.")
		return
	}
	notify := make(chan int, 1000)

	index := 0
	for {
		if index >= t {
			break
		}
		index = index + 1
		go func(tno int) {
			probe := NewProbe("tcp", notify)
			probe.Run(tno)
		}(index)
	}

	for {
		begin := time.Now()
		num := 0
		for {
			select {
			case <-notify:
				num = num + 1
			}
			end := time.Now()
			dur := end.Sub(begin)
			if dur >= 1*time.Minute {
				fmt.Println("sum:", num)
				break
			}
		}
	}
}
