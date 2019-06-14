package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/memberlist"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	bindAddr string
	bindPort string
	nodeName string
	cluster  string
)

func init() {
	flag.StringVar(&bindAddr, "ip", "127.0.0.1", "server ip")
	flag.StringVar(&bindPort, "port", "8000", "server port")
	flag.StringVar(&nodeName, "name", "test", "node name")
	flag.StringVar(&cluster, "c", "", "cluster address")
}

type simpleDelegate struct {
}

func (d *simpleDelegate) NodeMeta(limit int) []byte {
	return []byte("test")
}
func (d *simpleDelegate) NotifyMsg([]byte)                           {}
func (d *simpleDelegate) GetBroadcasts(overhead, limit int) [][]byte { return nil }
func (d *simpleDelegate) LocalState(join bool) []byte                { return nil }
func (d *simpleDelegate) MergeRemoteState(buf []byte, join bool)     {}

func main() {
	flag.Parse()
	eventCh := make(chan memberlist.NodeEvent, 10)

	config := memberlist.DefaultLANConfig()
	config.BindAddr = bindAddr
	config.BindPort, _ = strconv.Atoi(bindPort)
	config.Name = nodeName
	config.Events = &memberlist.ChannelEventDelegate{eventCh}
	config.GossipNodes = 1
	config.Delegate = &simpleDelegate{}
	config.LogOutput, _ = os.Create(os.DevNull)
	memberList, err := memberlist.Create(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		for {
			select {
			case e := <-eventCh:
				switch e.Event {
				case memberlist.NodeJoin:
					fmt.Println("name:", e.Node.Name)
					fmt.Println("addr:", e.Node.Addr.String())
					fmt.Println("port:", e.Node.Port)
					fmt.Println("meta:", e.Node.Meta)
				case memberlist.NodeLeave:
					fmt.Println(e.Node.Name)
				default:
					continue
				}
			}
		}
	}()

	_, err = memberList.Join(strings.Split(cluster, ","))
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, member := range memberList.Members() {
		fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
	}
	<-time.After(100000000 * time.Minute)
}
