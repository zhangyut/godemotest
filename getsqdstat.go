package main

import (
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
	"strings"
	"time"

	"quark"
	"quark/haproxy"
	"quark/registry"
	"sqd/statistics"
)

var (
	etcdcluster string
	local       string
)

func init() {
	flag.StringVar(&etcdcluster, "e", "http://localhost:2379", "etcd cluster address.")
	flag.StringVar(&local, "i", "127.0.0.1", "local ip addr")
}

func formatTime(t time.Time) string {
	return t.Format("2006.01.02 15:04:05")
}

func displayTable(data [][]string, headers []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetBorder(true)
	table.AppendBulk(data)
	table.Render()
	fmt.Printf("\n\n")
}

func displayStat(stat []statistics.StatResult) int {
	number := 0
	headers := []string{"type", "parent", "name", "count", "starttime"}
	bodys := [][]string{}
	for _, v := range stat {
		if v.Type == "topic" && len(bodys) > 0 {
			displayTable(bodys, headers)
			bodys = [][]string{}
			number++
		}

		prefix := ""
		if v.Type == "channel" {
			prefix = "--"
		} else if v.Type == "consumer" {
			prefix = "----"
		}
		tmp := []string{prefix + v.Type, v.Parent, v.Name, strconv.Itoa(int(v.Count)), formatTime(v.StartTime)}
		bodys = append(bodys, tmp)
	}
	displayTable(bodys, headers)
	number++
	return number
}

func main() {
	flag.Parse()
	registry, err := registry.NewEtcdRegistry(local, strings.Split(etcdcluster, ","))
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	//init nodeservice group and view
	proxy := haproxy.GetHttpCmdProxy(registry, "sqd_cmd", statistics.SupportedCmds())
	task := quark.NewTask()
	task.User = "admin"
	task.AddCmd(&statistics.Stat{
		Topic:   "",
		Channel: "",
	})
	ret := ""
	rets := []statistics.StatResult{}
	err = proxy.HandleTask(task, &rets, &ret)
	if err != nil {
		fmt.Println(err.Error())
	}
	if ret != "" {
		fmt.Println(ret)
	}
	number := displayStat(rets)
	fmt.Println("sum:" + strconv.Itoa(number) + " topics")
}
