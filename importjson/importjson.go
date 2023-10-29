package importjson

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// mongoDbへの接続情報構造体
type mongoImportParams struct {
	host           string
	user           string
	password       string
	databaseName   string
	collectionName string
}

// mongodbのパラメータを初期化する（新規構造体の作成）
func newMongoImportParams(host string, user string, password string, databaseName string, collectionName string) *mongoImportParams {
	mi := new(mongoImportParams)
	mi.host = host
	mi.user = user
	mi.password = password
	mi.databaseName = databaseName
	mi.collectionName = collectionName
	return mi
}

// jsonファイルをmongoDBにインポートする
func (mi mongoImportParams) importJson(filePath string) bool {
	// mongoimportを使用するためには、mongodb-database-toolsをインストールする必要がある。
	//https://www.mongodb.com/docs/database-tools/installation/installation-macos/
	cmd := exec.Command("mongoimport", "-h", mi.host, "-u", mi.user, "-p", mi.password, "--db", mi.databaseName, "--collection", mi.collectionName, "--file", filePath, "--jsonArray")
	_, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return true
}

// ファイルのパスを取得して配列にする。
func getFilePaths(dirPath string, extensionName string) []string {
	// ファイルのパスを格納する配列
	var filePaths []string

	// ディレクトリの中身を読み込む
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("error: path %v, err %v\n", path, err)
		}

		// ディレクトリは無視する
		if info.IsDir() {
			return nil
		}

		// 拡張子がにextensionNameであった場合は配列に格納する
		if strings.EqualFold(filepath.Ext(path), "."+extensionName) {
			filePaths = append(filePaths, path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %v: %v\n", dirPath, err)
	}

	return filePaths
}

func moveFile(src string, dstDir string) error {
	// ファイル名を取得する
	_, fileName := filepath.Split(src)

	// 移動先のディレクトリが存在しない場合は作成する
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		err := os.Mkdir(dstDir, 0777)
		if err != nil {
			return err
		}
	}
	// ファイルを移動する
	err := os.Rename(src, dstDir+"/"+fileName)
	if err != nil {
		return err
	}
	return err
}

//// SetDotenv .envファイルを読み込む
//func SetDotenv(envPath string) {
//	err := godotenv.Load(envPath)
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}
//}
//
//func init() {
//	log.SetFlags(log.Lshortfile)
//	// .envファイルを読み込む/
//	SetDotenv(".env")
//}

func example() {

	//import用の構造体を作成する
	mongoImport := newMongoImportParams(os.Getenv("MONGO_HOST"), os.Getenv("MONGO_USER"), os.Getenv("MONGO_PASSWORD"), os.Getenv("MONGO_DATABASE"), os.Getenv("MONGO_COLLECTION"))

	// 指定したディレクトリ内に存在するファイルから、指定した拡張子のファイルのパスを配列に格納する
	targetDirPath := "./input_data"
	choiceExtensionName := "json"
	jsonFilePath := getFilePaths(targetDirPath, choiceExtensionName)

	// jsonファイルをmongoDBにインポートする
	for _, filePath := range jsonFilePath {
		if mongoImport.importJson(filePath) {
			fmt.Println("import success")
			// インポートしたファイルを別のディレクトリに移動する
			err := moveFile(filePath, "./input_data/completed_data")
			if err != nil {
				return
			}

		}
	}
}
