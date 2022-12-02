package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	JSONMarshalExample()
	JSONUnmarshalExample()
}

type Address struct {
	Type    string
	City    string
	Country string
}

type VCard struct {
	FirstName string
	LastName  string
	Addresses []*Address
	Remark    string
}

func JSONMarshalExample() {
	a1 := &Address{
		Type:    "a",
		City:    "b",
		Country: "c",
	}
	vc := VCard{"Jan", "Kersschot", []*Address{a1, a1}, "none"}

	j, _ := json.Marshal(vc)
	fmt.Printf("JSON format: %s\n", j)

	// 写出到文件
	file, _ := os.OpenFile(".//vcard.json", os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	encoder := json.NewEncoder(file)
	err := encoder.Encode(vc)
	if err != nil {
		log.Println("Error in encoding json")
	}

}

type FamilyMember struct {
	Name    string
	Age     int
	Parents []string
}

func JSONUnmarshalExample() {
	b := []byte(`{"Name": "Wednesday", "Age": 6, "Parents": ["Gomez", "Morticia"]}`)
	var f any
	err := json.Unmarshal(b, &f)
	if err != nil {
		fmt.Println("json unmarshal error")
	}
	fmt.Println("f=\t", f)

	m := f.(map[string]any)
	fmt.Println("m=\t", m)

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", vv)

		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don’t know how to handle")
		}

	}

	var m1 FamilyMember
	//	解码到数据结构
	err = json.Unmarshal(b, &m1)
	if err != nil {
		fmt.Println("反序列化错误")
	}
	fmt.Println(m1)
}
