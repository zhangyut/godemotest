package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"strconv"
	"strings"
	"zcloud-go/rrmanager"
)

var (
	etcd    string
	localip string
	userid  string
)

func init() {
	flag.StringVar(&etcd, "e", "", "etcd")
	flag.StringVar(&localip, "l", "127.0.0.1", "local ip")
	flag.StringVar(&userid, "u", "", "user id")
}

func main() {
	flag.Parse()
	records := [][]string{
		{"域名ID(zcloud)", "租户名称", "租户ID", "应用基线", "配置费用", "域名", "区", "线路", "TTL", "类型", "RDATA", "区密码", "宕机切换", "探测方法", "tcp探测端口", "探测频率", "重试次数", "备份地址", "用户操作日志", "rr flag"},
	}

	if etcd == "" || userid == "" {
		fmt.Println("etcd or userid is nil.")
		return
	}
	register, err := registry.NewEtcdRegistry(localip, strings.Split(etcd, ","))
	if err != nil {
		fmt.Println("create registry failed:", err.Error())
		return
	}

	rrmanagerProxy := haproxy.GetRestProxy(register, "rrmanager", nil)
	zoneTask := quark.NewTask()
	zoneTask.User = userid
	zoneTask.AddCmd(&rest.GetCmd{
		ResourceType: "zone",
	})
	var ret string
	zones := []rrmanager.Zone{}
	err = rrmanagerProxy.HandleTask(zoneTask, &zones, &ret)
	if ret != "" {
		fmt.Println("get zones failed: " + ret)
	}
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, z := range zones {
		rrTask := quark.NewTask()
		rrTask.User = userid
		rrTask.AddCmd(&rest.GetCmd{
			ResourceType: "rr",
			Conds:        map[string]interface{}{"zone": z.Name},
		})
		var ret string
		rrs := []rrmanager.Rr{}
		err := rrmanagerProxy.HandleTask(rrTask, &rrs, &ret)
		if ret != "" {
			fmt.Println("get rr failed: " + ret)
			continue
		}
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		//	{"域名ID用于构造实例数据","租户名称", "租户ID", "应用基线", "配置费用", "域名", "区", "线路", "TTL", "类型", "RDATA", "区密码", "宕机切换", "探测方法", "tcp探测端口", "探测频率", "重试次数", "备份地址", "用户操作日志", "rr flag"}
		for _, r := range rrs {
			record := []string{}
			record = append(record, r.Id)                       //域名id
			record = append(record, "")                         //租户名称
			record = append(record, "")                         //租户id
			record = append(record, "")                         //应用基线
			record = append(record, "")                         //配置费用
			record = append(record, r.Name)                     //域名
			record = append(record, r.Zone)                     //区
			record = append(record, r.View)                     //线路
			record = append(record, strconv.Itoa(r.Ttl))        //TTL
			record = append(record, r.Type)                     //类型
			record = append(record, r.Rdata)                    //RDATA
			record = append(record, "")                         //区密码
			record = append(record, "false")                    //宕机切换
			record = append(record, "")                         //探测方法
			record = append(record, "")                         //tcp探测端口
			record = append(record, "")                         //探测频率
			record = append(record, "")                         //重试次数
			record = append(record, "")                         //备份地址
			record = append(record, "true")                     //用户操作日志
			record = append(record, strconv.Itoa(int(r.Flags))) //rr flag
			records = append(records, record)
		}
	}

	//w := csv.NewWriter(os.Stdout)
	fout, ioerr := os.Create("zclouddns_data.csv")
	if ioerr != nil {
		return
	}
	w := csv.NewWriter(fout)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	fout.Close()
}
