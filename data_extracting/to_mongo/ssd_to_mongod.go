package main

import (
	"context"
	"fmt"
	"github.com/akamensky/argparse"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go/types"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var OpenFolderError = "打开文件夹%v失败"
var ReadFileError = "读取文件%v错误"
var SerializationError = "%v序列化bson错误"
var ClipWriteError = "clip信息批量写入mongo出错 "

// var MetaWriteError = "Thread %d meta信息批量写入mongo出错: "
var MetaWriteError = "meta信息批量写入mongo出错: "
var MongoConfigError = "mongo Client配置错误: %s"
var MongoConnectionError = "mongo连接失败: %s"
var MongoCloseConnectionError = "mongodb连接关闭错误: %s"

func CheckError(err error, s string, a ...any) {
	if err != nil {
		log.Printf(strings.Join([]string{"[ERROR]", s}, "\t "), a...)
		log.Println("[ERROR]\t", err.Error())
	}
}

func CheckPanic(err error, s string, a ...any) {
	if err != nil {
		log.Printf(strings.Join([]string{"[ERROR]", s}, "\t"), a...)
		log.Println("[ERROR]\t", err.Error())
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

func (txn *TransactionExecutor) GetDocuments() error {
	ctx := txn.ctx
	client := txn.client
	db := txn.db
	collection := txn.coll
	coll := client.Database(db).Collection(collection)
	cursor, err := coll.Find(ctx, bson.D{})
	var results []bson.D
	if err = cursor.All(ctx, &results); err != nil {
		return err
	} else {
		for _, result := range results {
			fmt.Println(result)
		}
		return nil
	}
}

type TransactionExecutor struct {
	db       string        "dbname"
	coll     string        "coll name"
	client   *mongo.Client "mongo client object"
	bsonList []any         "bsonList"
	ctx      context.Context
	res      any "various operation result returned by mongo api"
	retryNum int "operation max retry num"
}

func (txn *TransactionExecutor) DoTransaction(transfunc func(sctx mongo.SessionContext) error, retryNum int) error {

	txn.retryNum = retryNum

	return txn.client.UseSessionWithOptions(
		txn.ctx, options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return txn.runTransactionWithRetry(sctx, transfunc)
		})

}

/*
*
sctx: 事务context
txnFn: 本次事务内要执行的操作

return err( TransientTransactionError(TTE), 返回事务执行错误, 如果是TTE就重试, 否则中断)
*/
func (txn *TransactionExecutor) runTransactionWithRetry(sctx mongo.SessionContext, txnFn func(sctx mongo.SessionContext) error) error {

	var err error
	for i := 0; i < txn.retryNum; i++ {
		err = txnFn(sctx) // Performs transaction.
		if err == nil {
			return nil
		}
		log.Println("Transaction aborted. Caught exception during transaction.")
		// If transient error, retry the whole transaction
		if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.HasErrorLabel("TransientTransactionError") {
			log.Println("TransientTransactionError, retrying transaction...")
			continue
		}
		return err
	}
	return err

}

/*
*
sctx session context
return error ( success if nil , retry if == UnknownTransactionCommitResult , should panic if other)
*/
func (txn *TransactionExecutor) commitWithRetry(sctx mongo.SessionContext) error {

	for {
		err := sctx.CommitTransaction(sctx)
		switch e := err.(type) {
		case nil:
			log.Println("Transaction committed.")
			return nil
		case mongo.CommandError:
			// Can retry commitf
			if e.HasErrorLabel("UnknownTransactionCommitResult") {
				log.Println("UnknownTransactionCommitResult, retrying commit operation...")
				continue
			}
			log.Println("Error during commit...")
			return e
		default:
			log.Println("Error during commit...")
			return e
		}
	}

}

//func MongoSessionExample() {
//	// simple example with no retry logic
//	// Create collections:
//	db.getSiblingDB("mydb1").foo.insertOne(
//	{abc: 0},
//	{ writeConcern: { w: "majority", wtimeout: 2000 } }
//	)
//	db.getSiblingDB("mydb2").bar.insertOne(
//	{xyz: 0},
//	{ writeConcern: { w: "majority", wtimeout: 2000 } }
//	)
//	// Start a session.
//	session = db.getMongo().startSession( { readPreference: { mode: "primary" } } );
//	coll1 = session.getDatabase("mydb1").foo;
//	coll2 = session.getDatabase("mydb2").bar;
//	// Start a transaction
//	session.startTransaction( { readConcern: { level: "local" }, writeConcern: { w: "majority" } } );
//	// Operations inside the transaction
//	try {
//		coll1.insertOne( { abc: 1 } );
//		coll2.insertOne( { xyz: 999 } );
//	} catch (error) {
//		// Abort transaction on error
//		session.abortTransaction();
//		throw error;
//	}
//	// Commit the transaction using write concern set at transaction start
//	session.commitTransaction();
//	session.endSession();
//}

func (txn *TransactionExecutor) TransactionInsertMany(sctx mongo.SessionContext) error {
	client := txn.client
	db := txn.db
	collection := txn.coll
	bsonList := txn.bsonList

	// Prereq: Create collections.
	coll := client.Database(db).Collection(collection)
	// configure write/read concern and read preference
	err := sctx.StartTransaction(options.Transaction().
		SetReadConcern(readconcern.Snapshot()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
	)
	if err != nil {
		return err
	}

	opts := options.InsertMany().SetOrdered(false).SetBypassDocumentValidation(false)
	//_, err = coll.InsertMany(sctx, []any{bson.M{"a":"b"}}, opts)
	res, err := coll.InsertMany(sctx, bsonList, opts)
	//fmt.Printf("err类型为 %T", err)
	if err != nil {
		sctx.AbortTransaction(sctx)
		log.Printf("caught exception during transaction while writing to %s/%s, aborting.", db, collection)
		return err
	}
	txn.res = res

	return txn.commitWithRetry(sctx)
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

// ...

//func test1() {
//
//	// ...
//	mongoURI := "mongodb://mongo.cowarobot.cn:27017"
//	database := "test"
//	clipColl := "test"
//
//	credential := options.Credential{
//		Username: "root",
//		Password: "cowa2022",
//	}
//
//	clientOpts := options.Client().ApplyURI(mongoURI).SetConnectTimeout(10 * 1e9).SetAuth(credential)
//
//	ctx := context.TODO()
//	client, err := mongo.Connect(ctx, clientOpts)
//
//	wc := writeconcern.New(writeconcern.WMajority())
//	rc := readconcern.Snapshot()
//	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
//
//	session, err := client.StartSession()
//	if err != nil {
//		panic(err)
//	}
//	defer session.EndSession(context.Background())
//
//	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
//		if err = session.StartTransaction(txnOpts); err != nil {
//			return err
//		}
//		result, err :=
//		if err != nil {
//			return err
//		}
//
//		if err != nil {
//			return err
//		}
//		if err = session.CommitTransaction(sessionContext); err != nil {
//			return err
//		}
//		return nil
//	})
//	if err != nil {
//		if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
//			panic(abortErr)
//		}
//		panic(err)
//	}
//}

func testTransaction() {
	// 连接信息初始化
	mongoURI := "mongodb://mongo.cowarobot.cn:27017"
	database := "test"
	clipColl := "test"
	ctx := context.TODO()
	credential := options.Credential{
		Username: "root",
		Password: "cowa2022",
	}
	// create client and ping
	clientOpts := options.Client().ApplyURI(mongoURI).SetConnectTimeout(10 * 1e9).SetAuth(credential)
	client, err := mongo.Connect(ctx, clientOpts)
	CheckPanic(err, MongoConfigError, mongoURI)
	// Ping the primary, test connection
	err = client.Ping(context.TODO(), readpref.Primary())
	CheckPanic(err, MongoConnectionError, mongoURI)
	// close connection
	defer func() {
		err := client.Disconnect(context.TODO())
		CheckPanic(err, MongoCloseConnectionError, mongoURI)
	}()

	bsonList := []any{bson.M{"a": "b"}, bson.M{"c": "d"}}
	texecutor := TransactionExecutor{
		db:       database,
		coll:     clipColl,
		client:   client,
		bsonList: bsonList,
		ctx:      ctx,
	}
	println("第一次获取, 应该什么都没有")
	err = texecutor.GetDocuments()
	CheckPanic(err, "获取document失败")

	err = texecutor.DoTransaction(texecutor.TransactionInsertMany, 3)
	CheckError(err, ClipWriteError)
	println("最后一次获取, 应该什么都没有")
	texecutor.GetDocuments()

}

func main() {
	//testTransaction()
	//
	//os.Exit(0)
	s1 := time.Now()
	metaWriteConcurrency, srcPath, mongoURI, database, clipColl, metaColl := argParser()
	ctx := context.TODO()

	credential := options.Credential{
		Username: "root",
		Password: "cowa2022",
	}
	// create client and ping
	clientOpts := options.Client().ApplyURI(mongoURI).SetConnectTimeout(10 * 1e9).SetAuth(credential)
	client, err := mongo.Connect(ctx, clientOpts)
	CheckPanic(err, MongoConfigError, mongoURI)
	// Ping the primary, test connection
	err = client.Ping(context.TODO(), readpref.Primary())
	CheckPanic(err, MongoConnectionError, mongoURI)

	// close connection
	defer func() {
		err := client.Disconnect(context.TODO())
		CheckPanic(err, MongoCloseConnectionError, mongoURI)
	}()

	// read files and turn into bsonList
	clipFiles, metaFiles := loadFile(srcPath)
	fmt.Printf("clips num, meta files num: %d,%d \n", len(clipFiles), len(metaFiles))

	clipBsonList := filesToBSONList(clipFiles)
	log.Printf("clips bson num: %d ", len(clipBsonList))
	println("concurrency:", metaWriteConcurrency)
	// insert clip data
	texecutor := TransactionExecutor{
		db:       database,
		coll:     clipColl,
		client:   client,
		bsonList: clipBsonList,
		ctx:      ctx,
	}
	err = texecutor.DoTransaction(texecutor.TransactionInsertMany, 3)
	switch v := texecutor.res.(type) {
	case *mongo.InsertManyResult:
		log.Printf("has written %d clip files to collection %s in database %s: \n", len(v.InsertedIDs), texecutor.coll, texecutor.db)
	case types.Nil:
		log.Printf("对库%v,%v集合的操作未生效", texecutor.db, texecutor.coll)
	}
	CheckPanic(err, ClipWriteError)

	metaBSONList := filesToBSONList(metaFiles)
	log.Printf("meta bson num: %d ", len(metaBSONList))
	te := TransactionExecutor{
		db:       database,
		coll:     metaColl,
		client:   client,
		bsonList: metaBSONList,
		ctx:      ctx,
	}
	err = te.DoTransaction(te.TransactionInsertMany, 3)
	switch v := te.res.(type) {
	case *mongo.InsertManyResult:
		log.Printf("has written %d meta files to collection %s in database %s: \n", len(v.InsertedIDs), te.coll, te.db)
	case types.Nil:
		log.Printf("对库%v,%v集合的操作未生效", te.db, te.coll)
	}
	CheckPanic(err, MetaWriteError)

	// insert meta data
	//wg.Add(metaWriteConcurrency)
	//for i := 0; i < metaWriteConcurrency; i++ {
	//	metaFilesSplit := metaFiles[i*len(metaFiles)/metaWriteConcurrency : (i+1)*len(metaFiles)/metaWriteConcurrency]
	//	go func(i int) {
	//		defer wg.Done()
	//		start := time.Now()
	//		metaBSONList := filesToBSONList(metaFilesSplit)
	//		te := TransactionExecutor{
	//			db:       database,
	//			coll:     metaColl,
	//			client:   client,
	//			bsonList: metaBSONList,
	//			ctx:      ctx,
	//		}
	//
	//		err = te.DoTransaction(te.TransactionInsertMany, 3)
	//		switch v := te.res.(type) {
	//		case *mongo.InsertManyResult:
	//			log.Printf("Thread %d has written %d meta files to collection %s in database %s: \n", i, len(v.InsertedIDs), te.coll, te.db)
	//		case types.Nil:
	//			log.Printf("对库%v,%v集合的操作未生效", te.db, te.coll)
	//		}
	//		CheckError(err, MetaWriteError, i)
	//		end := time.Now()
	//		log.Printf("Thread %d time cost: %s \n", i, end.Sub(start))
	//	}(i)
	//}
	//
	//wg.Wait()
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
