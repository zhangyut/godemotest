package main

import (
	"fmt"
)

func main() {
	zone := "zhanglei.test."
	qsZoneName := ""
	if string(zone[len(zone)-1]) == "." {
		qsZoneName = string(zone[0 : len(zone)-1])
	}

	fmt.Println(qsZoneName)
}
