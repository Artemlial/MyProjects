package main

import (
	// inner
	"context"
	"fmt"
	"log"
	"time"

	// outer
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Transaction struct {
	CCnum      string  `bson:"ccnum"`
	Date       string  `bson:"date"`
	Amount     float64 `bson:"amount"`
	Cvv        string  `bson:"cvv"`
	Expiration string  `bson:"exp"`
}

var client *mongo.Client
var collection *mongo.Collection

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	DB := client.Database("store")
	collection = DB.Collection("transactions")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	res := make([]Transaction, 0)
	if err = cursor.All(context.TODO(), &res); err != nil {
		log.Fatal(err)
	}

	for _, txn := range res {
		fmt.Println(txn.CCnum, txn.Date, txn.Amount, txn.Cvv, txn.Expiration)
	}
}
