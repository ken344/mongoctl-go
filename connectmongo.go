package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

// mongoDbへの接続情報構造体
type mongoParams struct {
	host           string
	user           string
	password       string
	databaseName   string
	collectionName string
	client         *mongo.Client
	database       *mongo.Database
	collection     *mongo.Collection
}

// ConnectClient mongodbとのclientを作成する
func (m *mongoParams) ConnectClient() *mongoParams {
	uri := fmt.Sprintf("mongodb://%s", m.host)
	// DB接続用のクレデンシャル情報を定義
	credential := options.Credential{
		Username: m.user,
		Password: m.password,
		//rootユーザー以外の場合は、そのユーザが管理されているDBを指定する。
		AuthSource: m.databaseName,
	}
	// DB接続用のオプションを定義.uriとクレデンシャル情報を指定する。
	clientOpts := options.Client().ApplyURI(uri).SetAuth(credential)
	// DB接続.クライアントを作成する。
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		panic(err)
	}
	// クライアントを構造体にセットし、返却する。
	m.client = client
	return m
}

// ConnectDatabase mongodbのdatabaseとの接続
func (m *mongoParams) ConnectDatabase() *mongoParams {
	// databaseを接続。
	m.database = m.client.Database(m.databaseName)
	//collectionを初期化（別のdbへ接続した場合を考慮し、collectionを初期化する）
	m.collection = nil
	return m
}

// ConnectCollection mongodbのcollectionとの接続
func (m *mongoParams) ConnectCollection() *mongoParams {
	// collectionを接続。
	m.collection = m.database.Collection(m.collectionName)
	return m
}

// mongodbのパラメータを初期化する（新規構造体の作成）
func newMongoParams(host string, user string, password string, databaseName string, collectionName string) *mongoParams {
	mg := new(mongoParams)
	mg.host = host
	mg.user = user
	mg.password = password
	mg.databaseName = databaseName
	mg.collectionName = collectionName
	return mg
}

// Disconnect mongodbとの接続を切断する
func (m *mongoParams) Disconnect() {
	err := m.client.Disconnect(context.Background())
	if err != nil {
		return
	}
}

// mongodbの操作
func (m *mongoParams) findOne(filter interface{}) *mongo.SingleResult {
	return m.collection.FindOne(context.Background(), filter)
}

func (m *mongoParams) findMultiple(filter interface{}) (*mongo.Cursor, error) {
	return m.collection.Find(context.Background(), filter)
}

func (m *mongoParams) insertOne(document interface{}) (*mongo.InsertOneResult, error) {
	return m.collection.InsertOne(context.Background(), document)
}

func (m *mongoParams) insertMany(documents []interface{}) (*mongo.InsertManyResult, error) {
	return m.collection.InsertMany(context.Background(), documents)
}

// SetDotenv .envファイルを読み込む
func SetDotenv(envPath string) {
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func init() {
	log.SetFlags(log.Lshortfile)
	// envファイルを読み込む
	SetDotenv(".env")
}

func main() {

	todofukenDb := newMongoParams(os.Getenv("MONGO_HOST"), os.Getenv("MONGO_USER"), os.Getenv("MONGO_PASSWORD"), os.Getenv("MONGO_DATABASE"), os.Getenv("MONGO_COLLECTION"))
	todofukenDb.ConnectClient().ConnectDatabase().ConnectCollection()
	defer todofukenDb.Disconnect()

	filter := bson.D{{Key: "test-key", Value: bson.D{{Key: "$exists", Value: true}}}}
	var results []bson.M
	cursor, err := todofukenDb.findMultiple(filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())
	if err = cursor.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}

}
