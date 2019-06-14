package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/pprof"
	"time"
)

var (
	url      string
	uri      string
	username string
	password string
	count    int
	begin    int
	disable  int
)

func init() {
	flag.StringVar(&url, "url", "https://127.0.0.1:4430", "zcloud url")
	flag.StringVar(&uri, "uri", "", "resource uri")
	flag.StringVar(&username, "u", "admin", "username")
	flag.StringVar(&password, "p", "password", "password")
	flag.IntVar(&count, "c", 100, "rr count")
	flag.IntVar(&begin, "b", 0, "rr count begin")
	flag.IntVar(&disable, "disable", 0, "diable check")
}

type ReqBody struct {
	Buff []byte
}

type Rr struct {
	Zone  string
	Name  string
	View  string
	Type  string
	Rdata string
	Ttl   int
}

func (p *Rr) toJson() string {
	return fmt.Sprintf("{\"zone\":\"%s\", \"name\":\"%s\", \"view\":\"%s\", \"type\":\"%s\", \"rdata\":\"%s\", \"ttl\":%d}", p.Zone, p.Name, p.View, p.Type, p.Rdata, p.Ttl)
}

type Rrs []Rr

func (p Rrs) toJson() string {
	tmp := ""
	for i, _ := range p {
		if i == 0 {
			tmp = p[i].toJson()
		} else {
			tmp = fmt.Sprintf("%s,%s", tmp, p[i].toJson())
		}
	}
	return fmt.Sprintf("[%s]", tmp)
}

func (p *ReqBody) Read(out []byte) (int, error) {
	if len(p.Buff) <= 0 {
		return 0, fmt.Errorf("no data")
	}
	out = p.Buff
	return len(p.Buff), nil
}

func (p *ReqBody) Write(in string) int {
	tmp := []byte(in)
	p.Buff = tmp
	return len(p.Buff)
}

func getToken(client *http.Client) (string, error) {
	resp, err := client.Post(fmt.Sprintf("%s/%s", url, "auth_cmd"), "application/json;charset=UTF-8", bytes.NewBuffer([]byte(fmt.Sprintf("{\"resource_type\":\"gen_api_token\", \"zdnsuser\":\"\", \"attrs\":{\"grant_type\":\"password\", \"username\":\"%s\", \"password\":\"%s\"}}", username, password))))
	if err != nil {
		return "", err
	}
	if resp != nil && resp.StatusCode != 200 {
		return "", fmt.Errorf("get token error:%s", resp.Status)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	type ApiToken struct {
		Token string
	}
	tmp := ApiToken{}

	err = json.Unmarshal(respBody, &tmp)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return tmp.Token, nil
}

func main() {
	// 根据命令行指定文件名创建 profile 文件
	f, err := os.Create("./addrrs.pprof")
	if err != nil {
		panic(err.Error())
	}
	// 开启 CPU profiling
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	flag.Parse()
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	token, err := getToken(client)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(token)
	rrs := Rrs{}
	for i := begin; i < count+begin; i++ {
		rr := Rr{
			Zone:  "zhanglei.test.",
			Name:  fmt.Sprintf("www%d.zhanglei.test.", i),
			View:  "others",
			Type:  "aw",
			Rdata: "10 1.1.1.1",
			Ttl:   600,
		}
		rrs = append(rrs, rr)
	}
	disable = 0
	if disable == 0 {
		resp, err := client.Post(fmt.Sprintf("%s/%s", url, uri), "application/json;charset=UTF-8", bytes.NewBuffer([]byte(fmt.Sprintf("{\"resource_type\":\"%s\", \"zdnsuser\":\"%s\", \"attrs\":%s}", "rr", token, rrs.toJson()))))
		if err != nil {
			fmt.Println(fmt.Sprintf("add rr failed %s", err.Error()))
		}
		if resp != nil && resp.StatusCode != 200 {
			fmt.Println(fmt.Sprintf("add rr failed %s", resp.Status))
		}
	} else {
		resp, err := client.Post(fmt.Sprintf("%s/%s", url, uri), "application/json;charset=UTF-8", bytes.NewBuffer([]byte(fmt.Sprintf("{\"resource_type\":\"%s\", \"disable_logic_check\":true, \"zdnsuser\":\"%s\", \"attrs\":%s}", "rr", token, rrs.toJson()))))
		if err != nil {
			fmt.Println(fmt.Sprintf("add rr failed %s", err.Error()))
		}
		if resp != nil && resp.StatusCode != 200 {
			fmt.Println(fmt.Sprintf("add rr failed %s", resp.Status))
		}
	}
}
