package main

import (
	"fmt"
	"quark/registry"
	"time"
)

func main() {
	reg, err := registry.NewEtcdClientEx("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		fmt.Println(err)
		return
	}

	/*
		err = reg.AddDir("/company")
		if err != nil {
			fmt.Println(err)
			return
		}
	*/
	err = reg.SetFilesWithExpire(map[string]string{"/company/name": "zdns", "/company/num": "150"}, 1*time.Minute)
	if err != nil {
		fmt.Println(err)
		return
	}

	reg1, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		fmt.Println(err)
		return
	}

	reg1.WatchStruct("/company/name", nil, func(key string, ret interface{}) { fmt.Println("delete key " + key) }, func(err string) { fmt.Println("err:" + err) })
	<-time.After(2 * time.Minute)
}
