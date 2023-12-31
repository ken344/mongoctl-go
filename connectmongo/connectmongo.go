package connectmongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

// MongoParams mongoDbへの接続情報構造体
type MongoParams struct {
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
func (m *MongoParams) ConnectClient() *MongoParams {
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
func (m *MongoParams) ConnectDatabase() *MongoParams {
	// databaseを接続。
	m.database = m.client.Database(m.databaseName)
	//collectionを初期化（別のdbへ接続した場合を考慮し、collectionを初期化する）
	m.collection = nil
	return m
}

// ConnectCollection mongodbのcollectionとの接続
func (m *MongoParams) ConnectCollection() *MongoParams {
	// collectionを接続。
	m.collection = m.database.Collection(m.collectionName)
	return m
}

// NewMongoParams mongodbのパラメータを初期化する（新規構造体の作成）
func NewMongoParams(host string, user string, password string, databaseName string, collectionName string) *MongoParams {
	mg := new(MongoParams)
	mg.host = host
	mg.user = user
	mg.password = password
	mg.databaseName = databaseName
	mg.collectionName = collectionName
	return mg
}

// Disconnect mongodbとの接続を切断する
func (m *MongoParams) Disconnect() {
	err := m.client.Disconnect(context.Background())
	if err != nil {
		return
	}
}

// FindOne dbから1つのドキュメントを取得する
// filter := bson.D{{"_id", id}}
func (m *MongoParams) FindOne(filter interface{}) *mongo.SingleResult {
	return m.collection.FindOne(context.Background(), filter)
}

// FindMultiple dbから複数のドキュメントを取得する
// filter := bson.D{{"name", "bob"}}
func (m *MongoParams) FindMultiple(filter interface{}) (*mongo.Cursor, error) {
	return m.collection.Find(context.Background(), filter)
}

// InsertOne dbに1つのドキュメントを挿入する
// document := bson.D{{"name", "pi"}, {"value", 3.14159}}
func (m *MongoParams) InsertOne(document interface{}) (*mongo.InsertOneResult, error) {
	return m.collection.InsertOne(context.Background(), document)
}

// InsertMany dbに複数のドキュメントを挿入する
// documents := []interface{}{
// bson.D{{"name", "Alice"}},
// bson.D{{"name", "Bob"}},
// }
func (m *MongoParams) InsertMany(documents []interface{}) (*mongo.InsertManyResult, error) {
	return m.collection.InsertMany(context.Background(), documents)
}

// FindKeyExists 指定のKeyが存在する（あるいは存在しない）ドキュメントを取得する
func (m *MongoParams) FindKeyExists(keyName string, isExists bool) (*mongo.Cursor, error) {
	return m.collection.Find(context.Background(), bson.D{{Key: keyName, Value: bson.D{{Key: "$exists", Value: isExists}}}})
}

// DeleteOne 指定のドキュメントを削除する
// filter := bson.D{{"name", "bob"}}
func (m *MongoParams) DeleteOne(filter interface{}) (*mongo.DeleteResult, error) {
	return m.collection.DeleteOne(context.Background(), filter)
}

// DeleteMany 指定のドキュメントを削除する
// filter := bson.D{{"name", "bob"}}
func (m *MongoParams) DeleteMany(filter interface{}) (*mongo.DeleteResult, error) {
	return m.collection.DeleteMany(context.Background(), filter)
}

// UpdateOne 指定のドキュメントを更新する
// filter := bson.D{{"_id", id}}
// update := bson.D{{"$set", bson.D{{"email", "newemail@example.com"}}}}
func (m *MongoParams) UpdateOne(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return m.collection.UpdateOne(context.Background(), filter, update)
}

// UpdateMany 指定のドキュメントを更新する
// filter := bson.D{{"birthday", today}}
// update := bson.D{{"$inc", bson.D{{"age", 1}}}}
func (m *MongoParams) UpdateMany(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return m.collection.UpdateMany(context.Background(), filter, update)
}

// UpdateByID 指定したIDのドキュメントを更新する
func (m *MongoParams) UpdateByID(id interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return m.collection.UpdateByID(context.Background(), id, update)
}

// ReplaceOne 指定のドキュメントを置き換える
// filter := bson.D{{"_id", id}}
// replacement := bson.D{{"location", "NYC"}}
func (m *MongoParams) ReplaceOne(filter interface{}, replacement interface{}) (*mongo.UpdateResult, error) {
	return m.collection.ReplaceOne(context.Background(), filter, replacement)
}

// FindOneAndDelete 指定のドキュメントを削除する
// filter := bson.D{{"_id", id}}
func (m *MongoParams) FindOneAndDelete(filter interface{}) *mongo.SingleResult {
	return m.collection.FindOneAndDelete(context.Background(), filter)
}

// FindOneAndReplace 指定のドキュメントを置き換える
// filter := bson.D{{"_id", id}}
// replacement := bson.D{{"location", "NYC"}}
// var replacedDocument bson.M
func (m *MongoParams) FindOneAndReplace(filter interface{}, replacement interface{}) *mongo.SingleResult {
	return m.collection.FindOneAndReplace(context.Background(), filter, replacement)
}

// FindOneAndUpdate 指定のドキュメントを更新する
// filter := bson.D{{"_id", id}}
// update := bson.D{{"$set", bson.D{{"email", "newemail@example.com"}}}}
// var updatedDocument bson.M
func (m *MongoParams) FindOneAndUpdate(filter interface{}, update interface{}) *mongo.SingleResult {
	return m.collection.FindOneAndUpdate(context.Background(), filter, update)
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
//	// envファイルを読み込む
//	SetDotenv(".env")
//}

func example() {

	todofukenDb := NewMongoParams(os.Getenv("MONGO_HOST"), os.Getenv("MONGO_USER"), os.Getenv("MONGO_PASSWORD"), os.Getenv("MONGO_DATABASE"), os.Getenv("MONGO_COLLECTION"))
	todofukenDb.ConnectClient().ConnectDatabase().ConnectCollection()
	defer todofukenDb.Disconnect()

	var results []bson.M
	var filter interface{}
	var cursor *mongo.Cursor
	var err error

	// キーに特定の値が含まれるドキュメントをすべて取得する)
	filter = bson.D{{"ja", "東京都"}}
	cursor, err = todofukenDb.FindMultiple(filter)
	if err != nil {
		log.Fatal(err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cursor, context.Background())

	if err = cursor.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
	// キーが存在するドキュメントを取得する($exists)
	cursor, err = todofukenDb.FindKeyExists("en", true)
	if err != nil {
		log.Fatal(err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, context.Background())

	if err = cursor.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}
