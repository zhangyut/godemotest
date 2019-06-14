package main

import (
	"fmt"
	"golang.org/x/net/webdav/xml"
)

var xmldata string = "<xml><ToUserName><![CDATA[gh_49f38d910277]]></ToUserName><FromUserName><![CDATA[ozK_9v0CLSEbl7KnCvX75SnEBr_Y]]></FromUserName><CreateTime>1486915221</CreateTime><MsgType><![CDATA[text]]></MsgType><Content><![CDATA[你好]]></Content><MsgId>6386252246558574548</MsgId></xml>"

type Message struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   xml.CharData
	FromUserName xml.CharData
	CreateTime   int64
	MsgType      xml.CharData
	Content      xml.CharData
	MsgId        xml.CharData
}

func main() {
	msg := Message{}
	err := xml.Unmarshal([]byte(xmldata), &msg)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}

	fmt.Println(string(msg.ToUserName))
	fmt.Println(string(msg.FromUserName))
	fmt.Println(msg.CreateTime)
	fmt.Println(string(msg.MsgType))
	fmt.Println(string(msg.Content))
	fmt.Println(string(msg.MsgId))
}
