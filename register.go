package main

import (
	"fmt"
	"path"
	"quark/registry"
)

func main() {
	c, err := registry.NewEtcdClientEx("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic(err.Error())
	}

	key := "/rec"
	tmp, err := c.GetDir(key)
	if err != nil {
		err = c.AddDir(key)
		if err != nil {
			fmt.Println(err.Error() + ":/rec")
		}
	} else {
		fmt.Println(tmp)
	}

	tmp, err = c.GetDir(path.Join(key, "beijing"))
	if err != nil {
		err = c.AddDir(path.Join(key, "beijing"))
		if err != nil {
			fmt.Println(err.Error() + ":" + path.Join(key, "beijing"))
		}
	} else {
		fmt.Println(tmp)
	}

	tmp, err = c.GetDir(path.Join(key, "beijing", "beijing1"))
	if err != nil {
		err = c.AddDir(path.Join(key, "beijing", "beijing1"))
		if err != nil {
			fmt.Println(err.Error() + ":" + path.Join(key, "beijing", "beijing1"))
		}
	} else {
		fmt.Println(tmp)
	}
	err = c.SetFiles(map[string]string{path.Join(key, "beijing", "beijing1", "ip"): "1.1.1.1"})
	if err != nil {
		fmt.Println(err.Error())
	}
}
