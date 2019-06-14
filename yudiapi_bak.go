package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	resp, _ := http.Get("https://api.damddos.com/v1")
	if resp != nil {
		res, _ := ioutil.ReadAll(resp.Body)
		if res != nil {
			fmt.Println(string(res))
		}
	}
}
