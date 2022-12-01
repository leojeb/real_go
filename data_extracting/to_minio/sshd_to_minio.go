package main

import (
	"context"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
	"strings"
)

func CheckError(err error, s string, a ...any) {
	if err != nil {
		log.Printf(strings.Join([]string{"[ERROR]", s}, "\t "), a...)
		log.Println("[ERROR]\t", err.Error())
	}
}

func loadFile(path string) (datafiles []string) {
	// 打开指定文件夹
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModeDir)
	CheckError(err, "打开文件夹%v失败", path)
	defer f.Close()
	// 读取目录下所有文件
	fileInfo, err := f.ReadDir(-1)
	DataFiles := make([]string, 0)

	for _, info := range fileInfo {
		filename := info.Name()
		DataFiles = append(DataFiles, filename)
	}

	return DataFiles
}

func argParser() string {

	parser := argparse.NewParser(
		"delete Minio Objects",
		"delete minio objects from corresponding local data path",
	)
	parser.SetHelp("h", "help")

	srcPath := parser.String("f", "folder",
		&argparse.Options{
			Required: true,
			Help:     "set src folder path, e.g. /src/",
		})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Print(parser.Usage(err))
		os.Exit(1)
	}
	return *srcPath
}

func getdbSize(minioClient *minio.Client) {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	objectCh := minioClient.ListObjects(ctx, "zgf-lidar", minio.ListObjectsOptions{
		Recursive: true,
	})
	i := 0
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}
		// fmt.Println(object)

		i++
		fmt.Println(i)
	}
	fmt.Println("文件个数为: ", i)
}

func main() {
	endpoint := "ossapi.cowarobot.cn:9000"
	accessKeyID := "cil2U1mEpfwVA5H7"
	secretAccessKey := "9BRdnsdSaa1IE392gkQLQk7HTA88T6zM"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now setup

	getdbSize(minioClient)
	// os.Exit(0)

	// folder := argParser()
	// fileNames := loadFile(folder)

	// opts := minio.RemoveObjectOptions {
	// 	GovernanceBypass: true,
	// }

	// for _, file:= range fileNames{

	// 	err = minioClient.RemoveObject(context.Background(), "zgf-lidar", file, opts)
	// 	if err != nil {
	// 		fmt.Println(file, " delete error", err)
	//         os.Exit(0)
	// 		return
	// 	}
	//     fmt.Println(file, "\tdeleted")
	// }

}
