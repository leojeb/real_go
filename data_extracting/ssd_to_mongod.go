package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CheckError(err error) {
	if err != nil {
		println("错误print输出:")
		fmt.Println(err.Error())
		log.Println(err.Error())
		panic(err)
	}
}

var mongo_uri = "mongodb://root:root@172.16.110.100:27017/?connectTimeoutMS=30000&maxPoolSize=100"

func loadFile(path string) []string {
	// 打开指定文件夹
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModeDir)
	if err != nil {
		log.Fatalln(err.Error())
		//os.Exit(0)
	}
	defer f.Close()
	// 读取目录下所有文件
	fileInfo, err := f.ReadDir(-1)
	fmt.Println(fileInfo)

	files := make([]string, 0)
	for _, info := range fileInfo {
		files = append(files, info.Name())
	}
	return files
}

func filesToBSONList(files []string, path string) []interface{} {

	var bsonList []interface{} = make([]interface{}, len(files))

	for i := 0; i < len(files); i++ {
		// 打开文件
		filepath := path + "\\" + files[i]
		f, err1 := os.Open(filepath)
		defer f.Close()
		CheckError(err1)

		// 读取文件
		r := bufio.NewReader(f)
		var buf []byte
		_, err := r.Read(buf)
		CheckError(err)

		// 转换成bson添加进数组
		bson, marshalErr := bson.Marshal(buf)
		CheckError(marshalErr)
		bsonList = append(bsonList, bson)
	}

	return bsonList
}

func BulkWrite(uri string, client mongo.Client, bsonList []interface{}) {
	coll := client.Database("admin").Collection("test-batch")
	//docs := []interface{}{
	//	bson.D{{"title", "Record of a Shriveled Datum"}, {"text", "No bytes, no problem. Just insert a document, in MongoDB"}},
	//	bson.D{{"title", "Showcasing a Blossoming Binary"}, {"text", "Binary data, safely stored with GridFS. Bucket the data"}},
	//}

	res, err := coll.InsertMany(context.TODO(), bsonList)
	if err != nil {
		panic(err)
	}
	fmt.Printf("inserted documents with IDs %v\n", res.InsertedIDs)
}

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	dir, _ := os.Getwd()
	println("当前路径", dir)
	files := loadFile("..\\inputs")

	fmt.Println(files)
	os.Exit(0)
	//bsonList := filesToBSONList(files, )
	//BulkWrite()

	coll := client.Database("admin").Collection("test-batch")
	title := "Back to Reality"

	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{"title", title}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", title)
		return
	}
	if err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
}
