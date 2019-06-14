package main

import (
	"fmt"
	"quark/rest"
	"time"
)

func main() {
	db, err := rest.OpenDB("127.0.0.1", "zdns", "zdns", "zdns")
	if err != nil {
		fmt.Println("open postgres failed:", err.Error())
		return
	}

	dbconnect_num := 0
	for {
		go func() {
			tx, err := db.Begin()
			if err != nil {
				fmt.Println("create tx failed: %v", err.Error())
				return
			}
			_, err = tx.PrepareAndExec("select count(*) from zc_rr")
			if err == nil {
				fmt.Println("select table succeed.")
			}

			<-time.After(100 * time.Second)
			tx.Commit()
		}()
		<-time.After(1 * time.Second)
		dbconnect_num = dbconnect_num + 1
		fmt.Println("connect number:", dbconnect_num)
	}
}
