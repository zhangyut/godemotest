package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {
	//curl -H 'API-Key: YJHWEPQQHXWBPIDIKOBS5HFA4WGKHGXG2U5Q' https://api.vultr.com/v1/dns/create_record --data 'domain=witcher.bid' --data 'name=www' --data 'type=A' --data 'data=202.173.9.25'
	url := "https://api.vultr.com/v1/dns/create_record"

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", url, strings.NewReader("domain=witcher.bid&name=www&type=A&data=202.173.9.29"))
	req.Header.Set("API-Key", "YJHWEPQQHXWBPIDIKOBS5HFA4WGKHGXG2U5Q")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	rbody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Println(string(rbody))
		return
	}
	fmt.Println(string(rbody))
}
