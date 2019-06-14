package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
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

	src := "KsjpOeQof9nLnmKlyAFVJIeMT9EPpLYm7MqOA3VyynCfipRZ9nvYQLKhFiyJ0HDbLsZ2Xw6eQZ627OBj9D3MvQgPS0ixM2iaZA1mzKAarj5x7wunkUXToiprfeK24DElPjn3V/tEjL5RmeffJjrNBxYpX/8rN0R4kveEgi6KMt0="
	bsrc, _ := base64.StdEncoding.DecodeString(src)

	rng1 := rand.Reader
	ori, err := rsa.DecryptPKCS1v15(rng1, priKey.(*rsa.PrivateKey), bsrc)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(ori))
}
