package main

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func main() {
	url := "mongodb://localhost:27017?ssl=false"

	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	// Use session as normal
	result := []bson.M{}
	c := session.DB("zdns").C("zcloud")
	err = c.Find(bson.M{}).All(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	session.Close()
}
