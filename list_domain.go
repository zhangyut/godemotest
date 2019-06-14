package main

import (
	"fmt"
	"net/http"
	//"encoding/json"
	"bytes"
	"io/ioutil"
)

const token = "7ca57a9f85a19a6e4b9a248c1daca185"
const url = "https://api.qssec.com"

func main() {
	//reqMap := map[string]string {
	//	"action" : "service_list",
	//	"token" : token,
	//}
	//data, _ := json.Marshal(reqMap)
	//buf := bytes.NewBuffer(data)
	str := fmt.Sprintf("action=domain_list&token=%s", token)
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
