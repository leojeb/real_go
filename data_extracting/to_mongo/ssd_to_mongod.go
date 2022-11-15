package main

import (
	"context"
	"fmt"
	"github.com/akamensky/argparse"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var wg sync.WaitGroup
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
		log.Printf(s, a...)
		log.Println("\t", err.Error())
	}
}

func CheckPanic(err error, s string, a any) {
	if err != nil {
		log.Printf(s, a)
		log.Println("\t", err.Error())
		panic(err)
	}
}

//TODO. 识别clips, npy.json, jpg.json等文件. 并行写入mongo

func loadFile(path string) (clip []string, meta []string) {
	// 打开指定文件夹
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModeDir)
	CheckError(err, "打开文件夹%v失败", path)
	defer f.Close()
	// 读取目录下所有文件
	fileInfo, err := f.ReadDir(-1)

	clipFiles := make([]string, 0)
	metaFiles := make([]string, 0)

	for _, info := range fileInfo {
		filename := info.Name()
		if strings.HasSuffix(filename, ".clip") {
			clipFiles = append(clipFiles, path+"/"+info.Name())
		} else if strings.HasSuffix(filename, ".json") {
			metaFiles = append(metaFiles, path+"/"+info.Name())
		}
	}
	return clipFiles, metaFiles
}

func filesToBSONList(files []string) []interface{} {

	var bsonList []interface{} = make([]interface{}, len(files))
	//println("files的个数", len(files))

	for i := 0; i < len(files); i++ {
		// 读取文件
		buf, err := os.ReadFile(files[i])
		CheckError(err, OpenFolderError, files[i])

		// 转换成bson添加进数组
		var res bson.M
		marshalErr := bson.UnmarshalExtJSON(buf, true, &res)
		CheckError(marshalErr, SerializationError, files[i])
		if marshalErr != nil {
			continue
		}
		// 添加ID字段
		pathSplits := strings.Split(files[i], "/")
		fileName := pathSplits[len(pathSplits)-1]
		docID := ""
		if strings.HasSuffix(fileName, ".json") {
			docID = strings.ReplaceAll(fileName, ".json", "")
		} else if strings.HasSuffix(fileName, ".clip") {
			docID = strings.ReplaceAll(fileName, ".clip", "")
		} else if docID == "" {
			log.Printf("文件名 %s 不符合规范", files[i])
		}
		res["_id"] = docID
		bsonList[i] = res
	}

	// 清洗掉nil值
	l := len(bsonList)
	for i := 0; i < l; i++ {
		if bsonList[i] == nil {
			bsonList = append(bsonList[:i], bsonList[i+1:]...)
			l = len(bsonList)
			i--
		}
	}
	return bsonList
}

func BulkWrite(database string, collection string, client *mongo.Client, bsonList []interface{}) (int, error) {

	//if len(bsonList) <= 0 {
	//	return 0, errors.New("传入列表为空")
	//}
	coll := client.Database(database).Collection(collection)
	opts := options.InsertMany().SetOrdered(false).SetBypassDocumentValidation(false)
	res, err := coll.InsertMany(context.TODO(), bsonList, opts)
	//CheckError(err, "插入mongodb错误")

	if res != nil {
		//fmt.Printf("inserted documents with IDs %v\n", res.InsertedIDs)
		return len(res.InsertedIDs), err
	} else {
		return 0, nil
	}

}

func argParser() (int, string, string, string, string, string) {

	parser := argparse.NewParser(
		"writeMongo",
		"bulk write clip and meta files into mongodb concurrently",
	)
	parser.SetHelp("h", "help")
	parallelism := parser.Int("p", "parallelism",
		&argparse.Options{
			Required: false,
			Help:     "set concurrency num, e.g. 10",
			Default:  20,
		})
	srcPath := parser.String("s", "source",
		&argparse.Options{
			Required: true,
			Help:     "set src folder path, e.g. /src/",
		})
	mongoURI := parser.String("u", "uri",
		&argparse.Options{
			Required: true,
			Help:     "set dest mongo uri",
		})
	dbName := parser.String("d", "database",
		&argparse.Options{
			Required: true,
			Help:     "set dest mongo dbName, e.g. cowa-3d",
		})
	clipCollName := parser.String("c", "clip-collection",
		&argparse.Options{
			Required: true,
			Help:     "set dest mongo clip collection name, e.g. clips",
		})
	metaCollName := parser.String("m", "meta-collection",
		&argparse.Options{
			Required: true,
			Help:     "set dest mongo meta collection name, e.g. metas",
		})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Print(parser.Usage(err))
		os.Exit(1)
	}

	//verbose := parser.Flag(
	//	"", "verbose", &argparse.Options{
	//		Help: "Verbose mode",
	//	},
	//)
	return *parallelism, *srcPath, *mongoURI, *dbName, *clipCollName, *metaCollName
}

func main() {

	metaWriteConcurrency, srcPath, mongoURI, database, clipColl, metaColl := argParser()

	s1 := time.Now()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI).SetConnectTimeout(10*1e9))
	CheckPanic(err, MongoConfigError, mongoURI)
	// Ping the primary, 连接测试
	err = client.Ping(context.TODO(), readpref.Primary())
	CheckPanic(err, MongoConnectionError, mongoURI)

	defer func() {
		err := client.Disconnect(context.TODO())
		CheckPanic(err, MongoCloseConnectionError, mongoURI)
	}()

	clipFiles, metaFiles := loadFile(srcPath)
	fmt.Printf("clips num, meta files num: %d,%d \n", len(clipFiles), len(metaFiles))
	clipBsonList := filesToBSONList(clipFiles)
	fileNum, err := BulkWrite(database, clipColl, client, clipBsonList)
	log.Println("clips written:", fileNum)
	CheckError(err, ClipWriteError)

	println("concurrency:", metaWriteConcurrency)
	wg.Add(metaWriteConcurrency)

	for i := 0; i < metaWriteConcurrency; i++ {
		metaFilesSplit := metaFiles[i*len(metaFiles)/metaWriteConcurrency : (i+1)*len(metaFiles)/metaWriteConcurrency]
		go func(i int) {
			defer wg.Done()
			start := time.Now()
			metaBSONList := filesToBSONList(metaFilesSplit)
			fileNum, err = BulkWrite(database, metaColl, client, metaBSONList)
			log.Printf("Thread %d has written %d meta files: \n", i, fileNum)
			CheckError(err, MetaWriteError, i)
			end := time.Now()
			log.Printf("Thread %d time cost: %s \n", i, end.Sub(start))
		}(i)
	}

	wg.Wait()
	log.Println("finished")
	e1 := time.Now()
	log.Println("total time cost: ", e1.Sub(s1))
	os.Exit(0)

	//coll := client.Database("admin").Collection("test-batch")
	//title := "Back to Reality"
	//
	//var result bson.M
	//err = coll.FindOne(context.TODO(), bson.D{{"title", title}}).Decode(&result)
	//if err == mongo.ErrNoDocuments {
	//	fmt.Printf("No document was found with the title %s\n", title)
	//	return
	//}
	//if err != nil {
	//	panic(err)
	//}
	//
	//jsonData, err := json.MarshalIndent(result, "", "    ")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%s\n", jsonData)

}
