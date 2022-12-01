package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"strings"
)

//	type MongoColletionUtils struct {
//		client *mongo.Client "mongo client对象"
//		// coll string "collection name"
//		// db string "db name"
//		// opts options "mongo options对象"
//	}
var OpenFolderError = "打开文件夹%v失败"
var ReadFileError = "读取文件%v错误"
var SerializationError = "%v序列化bson错误"
var ClipWriteError = "clip信息批量写入mongo出错 "
var MetaWriteError = "并发线程%d meta信息批量写入mongo出错: "
var MongoConfigError = "mongo Client配置错误: %s"
var MongoConnectionError = "mongo连接失败: %s"
var MongoCloseConnectionError = "mongodb连接关闭错误: %s"

func CheckError(err error, s string, a ...any) {
	if err != nil {
		log.Printf(strings.Join([]string{"[ERROR]", s}, "\t "), a...)
		log.Println("[ERROR]\t", err.Error())
	}
}

func CheckPanic(err error, s string, a any) {
	if err != nil {
		log.Printf(strings.Join([]string{"[ERROR]", s}, "\t"), a)
		log.Println("[ERROR]\t", err.Error())
		panic(err)
	}
}

/*
*
delete dirty data that is produced by failed/deleted data-extracting tasks according to their results preserved in ssds
*/
func DeleteDirtyData(client *mongo.Client, db string, collection string, taskName string) (int64, error) {

	coll := client.Database(db).Collection(collection)
	//opts := options.DeleteOptions
	filter := bson.D{{"from_bag", taskName}}
	res, err := coll.DeleteMany(context.TODO(), filter)
	if res != nil {
		//fmt.Printf("inserted documents with IDs %v\n", res.InsertedIDs)
		return res.DeletedCount, err
	} else {
		return 0, err
	}

}

func GetClips() {

}

func main() {
	mongoURI := "mongodb://root:cowa2022@172.16.100.107:27017"
	db := "zgf-lidar"
	clipColl := "clips"
	metaColl := "metas"

	credential := options.Credential{
		Username: "root",
		Password: "cowa2022",
	}
	clientOpts := options.Client().ApplyURI(mongoURI).SetConnectTimeout(10 * 1e9).SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	CheckPanic(err, MongoConfigError, mongoURI)

	// Ping the primary, 连接测试
	err = client.Ping(context.TODO(), readpref.Primary())
	CheckPanic(err, MongoConnectionError, mongoURI)
	// 关闭连接
	defer func() {
		err := client.Disconnect(context.TODO())
		if err != nil {
		}
		CheckPanic(err, MongoCloseConnectionError, mongoURI)
	}()

	// 删除clip错误数据
	count, err := DeleteDirtyData(client, db, clipColl, "20221116103528")
	CheckError(err, "clip删除失败")
	fmt.Println("删除clip数目", count)
	// 删除meta错误数据
	count, err = DeleteDirtyData(client, db, metaColl, "20221116103528")
	CheckError(err, "meta删除失败")
	fmt.Println("删除clip数目", count)

}
