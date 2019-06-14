package main

import (
	"encoding/json"
	"fmt"
	"quark"
	"quark/broadcast"
	"quark/rest"
	"zcloud-go/rrmanager"
)

type taskInJson struct {
	ResourceType      string            `json:"resource_type"`
	DisableLogicCheck bool              `json:"disable_logic_check"`
	Module            int               `json:"module"`
	AliasName         string            `json:"alias_name"`
	User              string            `json:"zdnsuser"`
	Attrs             []json.RawMessage `json:"attrs"`
}

func main() {
	sender, err := broadcast.NewKfkSender("127.0.0.1:9092", "test", "127.0.0.1")
	if err != nil {
		panic(err)
	}
	var resourceType string
	t := quark.NewTask()
	t.User = "admin"
	rr := &rrmanager.Rr{
		Id:    "111ddddddddddddd.zhanglei.test",
		Zone:  "zhanglei.test.",
		Name:  "www.zhanglei.test.",
		View:  "others",
		Type:  "a",
		Rdata: "127.0.0.3",
		Flags: 1,
		Ttl:   3600,
	}
	t.AddCmd(&rest.PostCmd{NewResource: rr})
	attrs := make([]json.RawMessage, len(t.Cmds), len(t.Cmds))
	for i, c_ := range t.Cmds {
		c, _ := c_.(*rest.PostCmd)
		if resourceType == "" {
			resourceType = string(rest.GetResourceType(c.NewResource))
		}
		attrs[i], _ = json.Marshal(c.NewResource)
	}
	var tt taskInJson
	tt.User = t.User
	tt.ResourceType = "rr"
	tt.Attrs = attrs
	tt.AliasName = "add"

	body, err := json.Marshal(tt)
	if err != nil {
		fmt.Println("marshal err:", err)
		return
	}

	err = sender.Send(body)
	if err != nil {
		fmt.Println("send err:", err)
	} else {
		fmt.Println("send succeed.")
	}
}
