package main

import (
	"fmt"
	"quark"
	"quark/group"
	"quark/registry"
	"time"
	node "zcloud-go/nodeservice"
)

var (
	supportedCmds = []quark.Command{
		&node.AddRr{},
		&node.DeleteRr{},
		&node.UpdateRr{},
		&node.AddAcl{},
		&node.UpdateAcl{},
		&node.DeleteAcl{},
		&node.AddView{},
		&node.UpdateView{},
		&node.DeleteView{},
		&node.AddZone{},
		&node.IgnoreView{},
		&node.DeleteZone{},
		&node.PauseGroup{},
		&node.ResumeGroup{},
		&node.GetGroupInfo{},
		&node.GetMemberInfo{},
		&node.UpdateFollowView{},
	}
)

func main() {
	registry, _ := registry.NewEtcdRegistry("192.168.79.20", []string{"http://202.173.9.103:2379"})
	_, err := group.NewSqGroup("maplenodes", "192.168.79.20", supportedCmds, registry, "./")
	if err != nil {
		fmt.Println(err.Error())
	}
	<-time.After(2 * time.Second)
}
