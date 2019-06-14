package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Header struct {
	Authorization string
	ContentType   string
	JsonRpcTonce  int64
}

func NewHeader() *Header {
	tonce := time.Now()
	micr := tonce.UnixNano() / int64(1000)
	return &Header{ContentType: "text/json", JsonRpcTonce: micr}
}

func (p *Header) genBodyHash(accessKey, secretKey []byte, body BodyNormal) []byte {
	buff := bytes.NewBuffer([]byte{})
	buff.Write([]byte("tonce=" + strconv.FormatInt(p.JsonRpcTonce, 10) + "&"))
	buff.Write([]byte("accesskey=" + string(accessKey) + "&"))
	buff.Write([]byte("method=" + body.Method + "&"))
	params, err := json.Marshal(body.Params)
	if err != nil {
		return nil
	} else {
		buff.Write([]byte("params=" + string(params)))
		mac := hmac.New(sha256.New, secretKey)
		mac.Write(buff.Bytes())
		expectedMAC := mac.Sum(nil)
		return expectedMAC
	}
}

func (p *Header) SetAuth(accessKey, secretKey []byte, body BodyNormal) error {
	hashBytes := p.genBodyHash(accessKey, secretKey, body)
	if hashBytes != nil {
	} else {
		return errors.New("generate hash failed.")
	}
	buff := bytes.NewBuffer(accessKey)
	buff.Write([]byte(":"))
	buff.Write(hashBytes)
	str := base64.StdEncoding.EncodeToString(buff.Bytes())
	p.Authorization = str
	return nil
}

func (p *Header) SetHeader(header *http.Header) {
	header.Add("Authorization", "Basic "+p.Authorization)
	header.Add("Content-Type", p.ContentType)
	header.Add("Json-Rpc-Tonce", strconv.FormatInt(p.JsonRpcTonce, 10))
}

type Body interface {
	Marshal() ([]byte, error)
	SetMethod(string)
	AddParam(interface{})
}

type BodyNormal struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

func NewBodyNormal() *BodyNormal {
	return &BodyNormal{Params: []interface{}{}}
}

func (p *BodyNormal) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *BodyNormal) SetMethod(method string) {
	p.Method = method
	p.Params = []interface{}{}
}
func (p *BodyNormal) AddParam(param interface{}) {
	p.Params = append(p.Params, param)
}

type OperationFlushName struct {
	Id           string      `json:"id"`
	IpAddr       string      `json:"ipAddr"`
	Method       string      `json:"method"`
	Fee          int         `json:"fee"`
	CreateTime   int64       `json:"createTime"`
	Source       int         `json:"source"`
	SourceObject interface{} `json:"sourceObject"`
	Status       int         `json:"status"`
	UnblockTime  int64       `json:"unblockTime"`
	Zone         int         `json:"zone"`
}

type FlushNameResultMsg struct {
	Opreation []OperationFlushName `json:"operation"`
	ErrorList []interface{}        `json:"errorList"`
}

type Result interface {
	Unmarshal(data []byte) error
}

type ResultFlushName struct {
	Result       FlushNameResultMsg `json:"result"`
	RemainAmount int                `json:"remainAmount"`
	Code         int                `json:"code"`
}

func (p *ResultFlushName) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

//get flush name

type Answer struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Ttl    int    `json:"ttl"`
	Record string `json:"record"`
}

type DNSNode struct {
	NodeIP  string   `json:"NodeIP"`
	Status  string   `json:"Status"`
	Answers []Answer `json:"Answer"`
}

type Record struct {
	Province string    `json:"Province"`
	DNSNodes []DNSNode `json:"DNSNodes"`
}

type GetFlushNameResultMsg struct {
	Status  string   `json:"status"`
	Records []Record `json:"records"`
}

type ResultGetFlushName struct {
	Result GetFlushNameResultMsg `json:"result"`
	Code   int                   `json:"code"`
}

func (p *ResultGetFlushName) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

type ErrorMsg struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type ResultError struct {
	Error ErrorMsg `json:"error"`
}

func (p *ResultError) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

var (
	ret  string = "{\"result\": {\"operation\": [{\"id\": \"100000\",\"ipAddr\": \"www.xyz.com\",\"method\": \"flushName\",\"fee\": 20,\"createTime\": 1444984146074,\"source\": 0,\"sourceObject\": null,\"status\": 0,\"unblockTime\": 0,\"zone\": 0}],\"errorList\": [ ]},\"remainAmount\": 99,\"code\": 200}"
	ret1 string = "{\"result\":{\"status\": \"success\",\"records\":[{\"Province\": \"京\", \"DNSNodes\": [{\"NodeIP\": \"118.118.118.1\",\"Status\": \"NOERROR\",\"Answer\": [{\"Name\": \"www.baidu.com.\",\"Type\": \"CNAME\",\"Ttl\": 511,\"Record\": \"www.a.shifen.com.\"}, {\"Name\": \"www.a.shifen.com.\",\"Type\": \"A\",\"Ttl\": 31,\"Record\": \"220.181.112.244\"}]},{\"NodeIP\": \"219.141.136.10\",\"Status\": \"NOERROR\",\"Answer\": [{\"Name\": \"www.baidu.com.\",\"Type\": \"CNAME\",\"Ttl\": 723,\"Record\": \"www.a.shifen.com.\"}, {\"Name\": \"www.a.shifen.com.\",\"Type\": \"A\",\"Ttl\": 344,\"Record\": \"220.181.111.188\"},{\"Name\": \"www.a.shifen.com.\",\"Type\": \"A\",\"Ttl\": 344,\"Record\": \"220.181.112.244\"}]},{\"NodeIP\": \"219.141.140.10\",\"Status\": \"NOERROR\",\"Answer\": [{\"Name\": \"www.baidu.com.\",\"Type\": \"CNAME\",\"Ttl\": 1134,\"Record\": \"www.a.shifen.com.\"},{\"Name\": \"www.a.shifen.com.\",\"Type\": \"A\",\"Ttl\": 481,\"Record\": \"220.181.111.188\"},{\"Name\": \"www.a.shifen.com.\",\"Type\": \"A\",\"Ttl\": 481,\"Record\": \"220.181.112.244\"}]}]}]},\"code\": 200}"
)

func main() {
	header := NewHeader()

	body := NewBodyNormal()
	body.SetMethod("flushName")
	body.AddParam("www.zhanglei.test.")
	bodyBuff, err := body.Marshal()
	if err != nil {
		panic(err.Error())
	}

	/*
		body.SetMethod("getFlushName")
		body.AddParam("100000000")
		str, err = MyFmt(body)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(str))
		}

		r1 := &ResultGetFlushName{}
		err = r1.Unmarshal([]byte(ret1))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(*r1)
		}
	*/

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	buff := bytes.NewBuffer(bodyBuff)
	req, err := http.NewRequest("POST", "https://api.damddos.com/v1", buff)
	req.Host = "api.demddos.com"
	err = header.SetAuth([]byte("sss"), []byte("kkkkkkkkk"), *body)
	if err != nil {
		panic(err.Error())
	}

	header.SetHeader(&req.Header)
	resp, _ := client.Do(req)
	if resp != nil {
		res, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(res))
		if res != nil {
			r := &ResultError{}
			err = r.Unmarshal([]byte(res))
			if err != nil {
				fmt.Println(err)
			} else if r.Error.Message == "" {
				r1 := &ResultFlushName{}
				err = r1.Unmarshal([]byte(res))
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(*r1)
				}
			} else {
				fmt.Println(*r)
			}
		}
	}
}
