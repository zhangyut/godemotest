package main

import (
	"crypto/rand"
	"crypto/rsa"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"time"
)

var (
	keyPath string
)

func init() {
	flag.StringVar(&keyPath, "k", "", "key path")
}

func main() {
	flag.Parse()
	if keyPath == "" {
		fmt.Println("param error")
		return
	}
	rng := rand.Reader
	//priKey, err := rsa.GenerateKey(rng, 256)
	rawKey := []byte{}
	rawKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	priKey, err := ssh.ParseRawPrivateKey(rawKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	pubKey := priKey.(*rsa.PrivateKey).Public()
	text := []byte("shi jie di yi deng.11111")
	data, err := rsa.EncryptPKCS1v15(rng, pubKey.(*rsa.PublicKey), text)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))

	<-time.After(2 * time.Second)
	rng1 := rand.Reader
	ori, err := rsa.DecryptPKCS1v15(rng1, priKey.(*rsa.PrivateKey), data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(ori))
}
