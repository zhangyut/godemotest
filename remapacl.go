package main

import (
	"flag"
	"fmt"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"strings"
	"zcloud-go/rrmanager"
)

var (
	mapr map[string]string = map[string]string{"europe": "z_europe",
		"cn_cu":             "y_unicom_china",
		"oceania":           "z_oceania",
		"cn_cm":             "y_cmcc_china",
		"c_cu_ningxia":      "y_unicom_ningxia",
		"c_cm_beijing":      "y_cmcc_beijing",
		"c_indonesia":       "z_indonesia",
		"c_ct_suzhou":       "y_telecom_jiangsu_suzhou",
		"c_ct_shandong":     "y_telecom_shandong",
		"c_cu_shanghai":     "y_unicom_shanghai",
		"c_cm_huanan":       "y_cmcc_guangdong,y_cmcc_guangxi,y_cmcc_hainan",
		"c_thailand":        "z_thailand",
		"c_ct_henan":        "y_telecom_henan",
		"c_cu_fujian":       "y_unicom_fujian",
		"c_cu_shandong":     "y_unicom_shandong",
		"c_cu_shanxi":       "y_unicom_shanxi",
		"c_cu_hainan":       "y_unicom_hainan",
		"c_cernet":          "y_cernet_china",
		"c_africa":          "z_afrika",
		"c_ct_jiangxi":      "y_telecom_jiangxi",
		"c_ct_hainan":       "y_telecom_hainan",
		"c_cu_jilin":        "y_unicom_jilin",
		"c_cu_hubei":        "y_unicom_hubei",
		"c_oceania":         "z_oceania",
		"c_ct_hunan":        "y_telecom_hunan",
		"c_north_america":   "z_north_america",
		"c_ct_guangdong":    "y_telecom_guangdong",
		"c_cu_jinan":        "y_unicom_shandong_jinan",
		"c_dp_huadong":      "y_drpeng_china_east",
		"c_dp_huanan":       "y_drpeng_china_southern",
		"c_ct_tianjign":     "y_telecom_tianjin",
		"c_telecom":         "y_telecom_china",
		"c_cu_hebei":        "y_unicom_hebei",
		"cn_ct":             "y_telecom_china",
		"sorth_america":     "z_south_america",
		"asia_other":        "",
		"cn_ck":             "y_drpeng_china",
		"north_america":     "z_north_america",
		"c_zealand":         "z_new_zealand",
		"c_dp_huazhong":     "y_drpeng_china_central",
		"c_cstnet":          "y_cstnet_china",
		"c_cu_heilongjiang": "y_unicom_heilj",
		"c_cm_huazhong":     "y_cmcc_henan,y_cmcc_hubei,y_cmcc_hunan",
		"c_cu_xinjiang":     "y_unicom_xinjiang",
		"c_australia":       "z_australia",
		"c_mauritius":       "z_mauritius",
		"c_ct_xinjiang":     "y_telecom_xinjiang",
		"c_ct_chongqing":    "y_telecom_chongqing",
		"c_cu_gansu":        "y_unicom_gansu",
		"c_cm_huadong":      "y_cmcc_shandong,y_cmcc_jiangsu,y_cmcc_anhui,y_cmcc_shanghai,y_cmcc_zhejiang,y_cmcc_jiangxi,y_cmcc_fujian",
		"c_canada":          "z_canada",
		"c_ct_jilin":        "y_telecom_jilin",
		"c_ct_dongguan":     "y_telecom_guangdong_dongguan",
		"c_cu_yunnan":       "y_unicom_yunnan",
		"c_cu_shaanxi":      "y_unicom_shaanxi",
		"c_french":          "z_france",
		"c_sorth_america":   "z_south_america",
		"c_ct_beijing":      "y_telecom_beijing",
		"c_italy":           "z_italy",
		"c_ct_zhejiang":     "y_telecom_zhejiang",
		"c_cu_zhejiang":     "y_unicom_zhejiang",
		"c_europe":          "z_europe",
		"c_cu_jiangxi":      "y_unicom_jiangxi",
		"c_ct_shanxi":       "y_telecom_shanxi",
		"c_ct_guangzhou":    "y_telecom_guangdong_guangzhou",
		"c_ct_shaanxi":      "y_telecom_shaanxi",
		"c_wasu":            "y_wasu_china",
		"africa":            "z_afrika",
		"cn_ws":             "y_wasu_china",
		"c_egypt":           "z_egypt",
		"c_cu_hunan":        "y_unicom_hunan",
		"c_ct_ningbo":       "y_telecom_zhejiang_ningbo",
		"c_cu_tianjing":     "y_unicom_tianjin",
		"c_holland":         "z_holland",
		"c_cuba":            "z_cuba",
		"c_cu_xizang":       "y_unicom_tibet",
		"c_drpeng":          "y_drpeng_china",
		"c_cu_guizhou":      "y_unicom_guizhou",
		"c_britain":         "z_britain",
		"c_ct_shenzhen":     "y_telecom_guangdong_shenzhen",
		"c_ct_ningxia":      "y_telecom_ningxia",
		"c_cu_jiangsu":      "y_unicom_jiangsu",
		"c_cu_guangxi":      "y_unicom_guangxi",
		"c_cn":              "z_china",
		"c_ct_fujian":       "y_telecom_fujian",
		"c_ct_shanghai":     "y_telecom_shanghai",
		"c_ct_hubei":        "y_telecom_hubei",
		"c_ct_sichuan":      "y_telecom_sichuan",
		"c_cu_neimeng":      "y_unicom_innermongolia",
		"c_cu_anhui":        "y_unicom_anhui",
		"c_cu_qingdao":      "y_unicom_shandong_qingdao",
		"c_cu_qinghai":      "y_unicom_qinghai",
		"c_ct_hebei":        "y_telecom_hebei",
		"c_ct_anhui":        "y_telecom_anhui",
		"c_ct_xizang":       "y_telecom_tibet",
		"c_ct_nanjing":      "y_telecom_jiangsu_nanjing",
		"c_cu_beijing":      "y_unicom_beijing",
		"c_cu_henan":        "y_unicom_henan",
		"c_cm_jiangsu":      "y_cmcc_jiangsu",
		"c_cm_xibei":        "y_cmcc_china_northwest",
		"c_singapore":       "z_singapore",
		"c_russia":          "z_russia",
		"c_cmcc":            "y_cmcc_china",
		"c_america":         "z_usa",
		"c_cu_guangdong":    "y_unicom_guangdong",
		"c_ct_guangxi":      "y_telecom_guangxi",
		"c_cu_liaoning":     "y_unicom_liaoning",
		"c_mexico":          "z_mexico",
		"c_ct_heilongjiang": "y_telecom_heilj",
		"c_cm_shandong":     "y_cmcc_shandong",
		"c_cm_dongbei":      "y_cmcc_china_northeast",
		"c_ct_jiangsu":      "y_telecom_jiangsu",
		"c_unicom":          "y_unicom_china",
		"c_seychelles":      "z_seychelles",
		"c_ct_hangzhou":     "y_telecom_zhejiang_hangzhou",
		"c_cu_chongqing":    "y_unicom_chongqing",
		"c_cm_huabei":       "y_cmcc_beijing,y_cmcc_tianjin,y_cmcc_hebei,y_cmcc_shanxi,y_cmcc_innermongolia",
		"c_ct_gansu":        "y_telecom_gansu",
		"c_ct_qinghai":      "y_telecom_qinghai",
		"c_ct_neimeng":      "y_telecom_innermongolia",
		"c_ct_zhongshan":    "y_telecom_guangdong_zhongshan",
		"c_korea":           "z_korea",
		"c_colombia":        "z_columbia",
		"c_ct_yunnan":       "y_telecom_yunnan",
		"c_cm_xinan":        "y_cmcc_sichuan,y_cmcc_guizhou,y_cmcc_yunnan,y_cmcc_tibet,y_cmcc_chongqing",
		"c_dp_huabei":       "y_drpeng_china_north",
		"c_cu_sichuan":      "y_unicom_sichuan",
		"c_japan":           "z_japan",
		"c_ct_liaoning":     "y_telecom_liaoning",
		"c_brazil":          "z_brazil",
		"c_ct_guizhou":      "y_telecom_guizhou",
		"c_malaysia":        "z_malaysia"}
	userId   string
	viewId   string
	local    string
	etcd     string
	onlyLook bool
)

func init() {
	flag.StringVar(&userId, "u", "", "user id")
	flag.StringVar(&viewId, "v", "", "view id")
	flag.StringVar(&local, "l", "127.0.0.1", "local ip")
	flag.StringVar(&etcd, "e", "http://127.0.0.1:2379", "etcd")
	flag.BoolVar(&onlyLook, "look", false, "only look")
}

func main() {
	flag.Parse()
	if userId == "" || viewId == "" {
		panic("param err.")
	}
	registry, err := registry.NewEtcdRegistry(local, strings.Split(etcd, ","))
	if err != nil {
		panic("create registry failed:" + err.Error())
	}
	store, err := rest.StoreForResources(rrmanager.SupportedResources(), registry)
	if err != nil {
		panic(err.Error())
	}
	tx, err := store.Begin()
	if err != nil {
		panic(err.Error())
	}
	for k, v := range mapr {
		acls := []rrmanager.Acl{}
		err := tx.Fill(map[string]interface{}{"id": k}, &acls)
		if err != nil {
			fmt.Println(err.Error())
		}
		if len(acls) != 1 {
			fmt.Println("acl error:" + k)
		}
		for _, m := range strings.Split(v, ",") {
			acls := []rrmanager.Acl{}
			err := tx.Fill(map[string]interface{}{"id": m}, &acls)
			if err != nil {
				fmt.Println(err.Error())
			}
			if len(acls) != 1 {
				fmt.Println("2 acl error:" + m)
			}
		}
	}
	views := []rrmanager.View{}
	err = tx.Fill(map[string]interface{}{"id": viewId}, &views)
	if len(views) == 0 {
		fmt.Println("no view:" + viewId)
		return
	}
	view := views[0]
	fmt.Println(view.Acls)
	newAcls := []string{}
	for _, m := range view.Acls {
		tmp, ok := mapr[m]
		if ok {
			newAcls = append(newAcls, strings.Split(tmp, ",")...)
		}
	}
	fmt.Println(newAcls)
	tx.Commit()
	fmt.Println(onlyLook)
	if onlyLook {
		return
	}

	view.Acls = newAcls
	proxy := haproxy.GetRestProxy(registry, "rrmanager", rrmanager.SupportedResources())

	t := quark.NewTask()
	t.User = userId
	t.AddCmd(&rest.PutCmd{NewResource: &view})
	var errInfo string
	err = proxy.HandleTask(t, nil, &errInfo)
	if err != nil {
		panic(err.Error())
	}
	if errInfo != "" {
		panic(errInfo)
	}
}
