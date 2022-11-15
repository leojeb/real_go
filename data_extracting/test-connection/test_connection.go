package main

import (
	"errors"
	"log"
)

const uri = "mongodb://root:root@172.16.110.100:27017"

func CheckError(err error, s string, a ...any) {
	if err != nil {
		log.Printf(s, a...)
		log.Println("\t", err.Error())
	}
}
func main() {
	// Create a new client and connect to the server
	//client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetConnectTimeout(10 * 1e9))
	//if err != nil {
	//	panic(err)
	//}
	//defer func() {
	//	if err = client.Disconnect(context.TODO()); err != nil {
	//		panic(err)
	//	}
	//}()
	//// Ping the primary
	//if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
	//	panic(err)
	//}
	//fmt.Println("Successfully connected and pinged.")

	CheckError(errors.New("测试"), "asdlifaf")
}
