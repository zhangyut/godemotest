package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	type RGB struct {
		R string
		G string
		B string
		A int
	}

	var j = []byte(`{"R": "RGB","G":"indid","B":"laidnnfi", "a":10}`)
	var colors RGB
	err := json.Unmarshal(j, &colors)
	if err != nil {
		log.Fatalln("error:", err)
	}
	fmt.Println(colors)

	type AddUser struct {
		Name             string `json:"name"`
		Telephone        string `json:"telephone"`
		Email            string `json:"email"`
		Company          string `json:"company"`
		MaxZoneCount     int    `json:"max_zone_count"`
		MaxRrCount       int    `json:"max_rr_count"`
		MaxView          int    `json:"max_view"`
		MaxAcl           int    `json:"max_acl"`
		MaxZdnsuserView  int    `json:"max_zdnsuser_view"`
		MaxAwCount       int    `json:"max_aw_count"`
		MaxCnamewCount   int    `json:"max_cnamew_count"`
		MaxXwCount       int    `json:"max_xw_count"`
		MaxFailoverCount int    `json:"max_failover_count"`
		MaxRedirectCount int    `json:"max_redirect_count"`
		MaxShortUrlCount int    `json:"max_short_url_count"`
		EnableSms        bool   `json:"enable_sms"`
	}

	var test = []byte(`{"email":"zhanglei@zdns.cn","company":"zdns","max_zone_count":100,"max_rr_count":100,"max_view":100,"max_acl":100,"max_zdnsuser_view":100,"max_aw_count":100,"max_cnamew_count":100,"max_xw_count":100,"max_failover_count":100,"max_rediret_count":100,"max_short_url_count":100,"enable_sms":true}`)
	var adduser AddUser
	adduser.Name = "empty"
	adduser.Telephone = "empty"
	err = json.Unmarshal(test, &adduser)
	if err != nil {
		log.Fatalln("error:", err)
	}
	fmt.Println(adduser)
}
