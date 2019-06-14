package main

import "reflect"
import "io"
import "fmt"

type T struct {
	A int `sql:"referto, uk"`
	B string
}

type MyReader struct{ Name string }

func (r MyReader) Read(p []byte) (n int, err error) {
	return 0, nil
}

func main() {
	t := T{203, "mh203"}
	tmp := reflect.ValueOf(t)
	tmpT := reflect.TypeOf(t)
	s := reflect.ValueOf(&t).Elem()
	typeOfT := s.Type()
	fmt.Printf("type name: %s, type kind: %v,  value name: %s, value kind: %v\n", typeOfT.Name(), tmpT.Kind(), s.Type().Name(), tmp.Kind())
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s %s = %v\n", i,
			typeOfT.Field(i).Name, typeOfT.Field(i).Type, f.Type(), f.Interface())
	}

	var reader io.Reader
	reader = &MyReader{"a.txt"}
	tmp1V := reflect.ValueOf(reader)
	tmp1T := reflect.TypeOf(reader)
	fmt.Printf("type name: %v, value name: %v\n", tmp1T.Name(), tmp1V.Elem().Type().Name())
	for i := 0; i < typeOfT.NumField(); i++ {
		field := typeOfT.Field(i)
		fieldTag := field.Tag.Get("sql")
		fmt.Println(fieldTag)
	}
}
