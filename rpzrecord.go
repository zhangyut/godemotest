package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"text/template"
)

const rpzconf = `$TTL 1H
@	SOA LOCALHOST. named-mgr.example.com (1 1h 15m 30d 2h)
	NS  LOCALHOST.
`

const rpzrecord = `{{.Name}}		A		{{.Value}}
`

type RpzRecord struct {
	Name  string
	Value string
}

func main() {
	var err error
	fout, err := os.Create(path.Join(".", "rpzconf.bak"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer fout.Close()
	bw := bufio.NewWriter(fout)

	t := template.Must(template.New("rpzconf").Parse(rpzconf))
	err = t.Execute(bw, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	rpzRecords := []RpzRecord{
		RpzRecord{"zhanglei.test", "1.1.1.1"},
		RpzRecord{"zhanglei1.test", "2.1.1.1"},
		RpzRecord{"zhanglei2.test", "3.1.1.1"},
	}
	for _, v := range rpzRecords {
		tmp := v
		t := template.Must(template.New("rpzrecord").Parse(rpzrecord))
		err = t.Execute(bw, tmp)
		if err == nil {
			os.Rename("rpzconf.bak", "rpzconf")
		} else {
			fmt.Println(err.Error())
			return
		}
	}

	bw.Flush()
	return
}
