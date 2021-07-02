package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func Test(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	database := client.Database("gameOfLife")
	if err := database.Drop(ctx); err != nil {
		panic(err)
	}

	collection := database.Collection("gameOfLife")

	if _, err := collection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    &bson.D{{"coordinate.0", 1}, {"coordinate.1", 1}},
			Options: (&options.IndexOptions{}).SetUnique(true),
		},
	); err != nil {
		panic(err)
	}

	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			cellData := CellData{
				Coordinate: [2]int{x, y},
				Cell:       CellLive,
			}
			if _, err = collection.InsertOne(ctx, &cellData); err != nil {
				panic(err)
			}
		}
	}

	cur, err := collection.Aggregate(ctx, mongo.Pipeline{
		{{"$match", bson.D{{"$or", bson.A{
			bson.D{{"coordinate", bson.A{0, 0}}},
			bson.D{{"coordinate", bson.A{1, 0}}},
			bson.D{{"coordinate", bson.A{2, 0}}},
		}}}}},
		{{"$group", bson.D{{"_id", ""}, {"sum", bson.D{{"$sum", "$cell"}}}}}},
	})
	if err != nil {
		panic(err)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result bson.M
		if err := cur.Decode(&result); err != nil {
			panic(err)
		}
		fmt.Println(result)
	}

	if cur.Err() != nil {
		panic(cur.Err())
	}
}
