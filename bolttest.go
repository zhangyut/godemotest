package main

import (
	"dragondance/backend/bolt"
	"flag"
	"fmt"
	"sort"
)

type SortedArray struct {
	Ids []string
}

func (p *SortedArray) Len() int {
	return len(p.Ids)
}
func (p *SortedArray) Less(i, j int) bool {
	return p.Ids[i] < p.Ids[j]
}
func (p *SortedArray) Swap(i, j int) {
	tmp := p.Ids[i]
	p.Ids[i] = p.Ids[j]
	p.Ids[j] = tmp
}

func main() {
	diff := SortedArray{Ids: []string{}}
	db52 := SortedArray{Ids: []string{}}
	db53 := SortedArray{Ids: []string{}}
	db59 := SortedArray{Ids: []string{}}

	var dbname string
	flag.StringVar(&dbname, "d", "", "db name")
	flag.Parse()
	if dbname == "" {
		//fmt.Println("db name can't nil.")
		//return
	}
	dbwrapper, err := bolt.NewBoltKVStore("bolt52.db")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	db, err := dbwrapper.Select("failover-rrset")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	myPrint := func(key string, value []byte) error {
		db52.Ids = append(db52.Ids, key)
		return nil
	}
	sort.Sort(&db52)
	err = db.ForEachKeyValue(myPrint)
	if err != nil {
		fmt.Println("read db52 err:", err.Error())
		return
	}

	///////////////
	dbwrapper, err = bolt.NewBoltKVStore("bolt53.db")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	db, err = dbwrapper.Select("failover-rrset")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	myPrint = func(key string, value []byte) error {
		db53.Ids = append(db53.Ids, key)
		return nil
	}
	sort.Sort(&db53)
	err = db.ForEachKeyValue(myPrint)
	if err != nil {
		fmt.Println("read db53 err:", err.Error())
		return
	}
	///////////////
	dbwrapper, err = bolt.NewBoltKVStore("bolt59.db")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	db, err = dbwrapper.Select("failover-rrset")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	myPrint = func(key string, value []byte) error {
		db59.Ids = append(db59.Ids, key)
		return nil
	}
	sort.Sort(&db59)
	err = db.ForEachKeyValue(myPrint)
	if err != nil {
		fmt.Println("read db59 err:", err.Error())
		return
	}

	for i, v := range db52.Ids {
		index := sort.SearchStrings(db53.Ids, v)
		if index == len(db53.Ids) || db53.Ids[index] != v {
			diff.Ids = append(diff.Ids, db52.Ids[i])
		}
	}

	for i, v := range db52.Ids {
		index := sort.SearchStrings(db59.Ids, v)
		if index == len(db59.Ids) || db59.Ids[index] != v {
			diff.Ids = append(diff.Ids, db52.Ids[i])
		}
	}

	for i, v := range db53.Ids {
		index := sort.SearchStrings(db52.Ids, v)
		if index == len(db52.Ids) || db52.Ids[index] != v {
			diff.Ids = append(diff.Ids, db53.Ids[i])
		}
	}

	for i, v := range db53.Ids {
		index := sort.SearchStrings(db59.Ids, v)
		if index == len(db59.Ids) || db59.Ids[index] != v {
			diff.Ids = append(diff.Ids, db53.Ids[i])
		}
	}

	for i, v := range db59.Ids {
		index := sort.SearchStrings(db52.Ids, v)
		if index == len(db52.Ids) || db52.Ids[index] != v {
			diff.Ids = append(diff.Ids, db59.Ids[i])
		}
	}

	for i, v := range db59.Ids {
		index := sort.SearchStrings(db53.Ids, v)
		if index == len(db53.Ids) || db53.Ids[index] != v {
			diff.Ids = append(diff.Ids, db59.Ids[i])
		}
	}

	sort.Sort(&diff)
	diffArray := []string{}
	for i, v := range diff.Ids {
		if i != 0 {
			if diffArray[len(diffArray)-1] == v {
				continue
			} else {
				diffArray = append(diffArray, diff.Ids[i])
			}
		} else {
			diffArray = append(diffArray, diff.Ids[i])
		}
	}

	for _, v := range diffArray {
		fmt.Println(v)
	}
}
