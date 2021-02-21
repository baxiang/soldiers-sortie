package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Student struct {
	ID     string `bson:"_id,omitempty"`
	Name   string `bson:"name"` // 名字
	Gender string `bson:gender` // 性别
	Age    int8   `bson:"age"`  // 年龄
}

func main() {
	clientOptions := options.Client().SetHosts([]string{"localhost:27017"}).
		SetConnectTimeout(5 * time.Second).
		SetAuth(options.Credential{Username: "admin", Password: "654321"})
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, clientOptions)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("mongodb connect success")

	collection := client.Database("my_db").Collection("my_collection")
	filter := bson.D{{"name", "张三"}}
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("delete count", res.DeletedCount)
}
