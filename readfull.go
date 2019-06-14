package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	var msgSize int32
	curFileName := "./data.txt"
	readFile, err := os.OpenFile(curFileName, os.O_RDONLY, 0600)
	if err != nil {
		panic(err.Error())
	}

	<-time.After(20 * time.Second)
	reader := bufio.NewReader(readFile)

	var totle int32
	for {
		fmt.Println("totle:", totle)
		err = binary.Read(reader, binary.BigEndian, &msgSize)
		if err != nil {
			readFile.Close()
			readFile = nil
			panic(err.Error())
		}

		readBuf := make([]byte, msgSize)
		_, err := io.ReadFull(reader, readBuf)
		if err != nil {
			readFile.Close()
			readFile = nil
			panic(err.Error())
		}
		totle = totle + int32(4) + msgSize
	}
}
