package main

import (
	"fmt"
	"net/http"
	//"encoding/json"
	"bytes"
	"flag"
	"io/ioutil"
)

const token = "7ca57a9f85a19a6e4b9a248c1daca185"
const url = "https://api.qssec.com"

func main() {
	var domainId string
	flag.StringVar(&domainId, "d", "", "domain id")
	flag.Parse()
	if domainId == "" {
		fmt.Println("domain id can't nil")
		return
	}
	//reqMap := map[string]string {
	//	"action" : "service_list",
	//	"token" : token,
	//}
	//data, _ := json.Marshal(reqMap)
	//buf := bytes.NewBuffer(data)
	str := fmt.Sprintf("action=record_list&domain_id=%s&token=%s", domainId, token)
	buf := bytes.NewBuffer([]byte(str))
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		fmt.Printf("NewRequest failed:%s.\n", err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("client do failed:%s.\n", err.Error())
		return
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
