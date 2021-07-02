package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"time"
)

type DatabaseBoard struct {
	Collection mongo.Collection
	Client     mongo.Client
}

const collectionName = "gameOfLife"

type CellData struct {
	ID         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Coordinate [2]int
	Cell
}

var _ Board = &DatabaseBoard{}

func newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second)
}

func NewDatabaseBoard() (DatabaseBoard, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return DatabaseBoard{}, err
	}
	{
		ctx, cancel := newContext()
		defer cancel()

		if err := client.Connect(ctx); err != nil {
			return DatabaseBoard{}, err
		}

		if err := client.Ping(ctx, nil); err != nil {
			return DatabaseBoard{}, err
		}
	}
	var collection *mongo.Collection
	{
		ctx, cancel := newContext()
		defer cancel()

		database := client.Database("gameOfLife")
		if err := database.Drop(ctx); err != nil {
			return DatabaseBoard{}, err
		}

		collection = database.Collection(collectionName)

		if _, err := collection.Indexes().CreateOne(
			ctx,
			mongo.IndexModel{
				Keys:    &bson.D{{"coordinate.0", 1}, {"coordinate.1", 1}},
				Options: (&options.IndexOptions{}).SetUnique(true),
			},
		); err != nil {
			return DatabaseBoard{}, err
		}
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	for x := 0; x < BoardWidth; x++ {
		for y := 0; y < BoardHeight; y++ {
			if err := insertCell(collection, x, y, Cell(random.Int63()&1)); err != nil {
				return DatabaseBoard{}, err
			}
		}
	}

	return DatabaseBoard{
		Collection: *collection,
	}, nil
}

func insertCell(collection *mongo.Collection, x int, y int, cell Cell) error {
	ctx, cancel := newContext()
	defer cancel()

	cellData := CellData{
		ID:         primitive.NewObjectID(),
		Coordinate: [2]int{x, y},
		Cell:       cell,
	}

	_, err := collection.InsertOne(ctx, &cellData)
	return err
}

func (d *DatabaseBoard) Close() error {
	ctx, cancel := newContext()
	defer cancel()
	return d.Client.Disconnect(ctx)
}

func (d *DatabaseBoard) Cell(x int, y int) (Cell, error) {
	ctx, cancel := newContext()
	defer cancel()
	result := d.Collection.FindOne(ctx, bson.M{"coordinate.0": x, "coordinate.1": y})
	if result.Err() != nil {
		return CellDead, result.Err()
	}

	cellData := &CellData{}
	return cellData.Cell, result.Decode(cellData)
}

func (d *DatabaseBoard) Set(x int, y int, cell Cell) error {
	ctx, cancel := newContext()
	defer cancel()
	_, err := d.Collection.UpdateOne(
		ctx,
		bson.M{"coordinate.0": x, "coordinate.1": y},
		mongo.Pipeline{{{"$set", bson.M{"cell": cell}}}},
	)
	return err
}

func wrapAround(a int, b int, max int) int {
	sum := a + b
	if sum < 0 {
		return max + sum
	}
	return sum % max
}

func coordinates(x int, y int) bson.A {
	leftColumn := wrapAround(x, -1, BoardWidth)
	rightColumn := wrapAround(x, 1, BoardWidth)
	topRow := wrapAround(y, -1, BoardHeight)
	bottomRow := wrapAround(y, 1, BoardHeight)
	return bson.A{
		bson.D{
			{"coordinate.0", leftColumn},
			{"coordinate.1", topRow},
		},
		bson.D{
			{"coordinate.0", x},
			{"coordinate.1", topRow},
		},
		bson.D{
			{"coordinate.0", rightColumn},
			{"coordinate.1", topRow},
		},
		bson.D{
			{"coordinate.0", leftColumn},
			{"coordinate.1", y},
		},
		bson.D{
			{"coordinate.0", rightColumn},
			{"coordinate.1", y},
		},
		bson.D{
			{"coordinate.0", leftColumn},
			{"coordinate.1", bottomRow},
		},
		bson.D{
			{"coordinate.0", x},
			{"coordinate.1", bottomRow},
		},
		bson.D{
			{"coordinate.0", rightColumn},
			{"coordinate.1", bottomRow},
		},
	}
}

func getSum(collection *mongo.Collection, x int, y int) (int, error) {
	ctx, cancel := newContext()
	defer cancel()

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{
		{{
			"$match", bson.D{{"$or", coordinates(x, y)}},
		}},
		{{"$group", bson.D{{"_id", "sum"}, {"sum", bson.D{{"$sum", "$cell"}}}}}},
	})
	if err != nil {
		return 0, err
	}

	if cursor.Next(ctx) == false {
		return 0, errors.New("0 rows returned")
	}

	sum := &struct{ Sum int }{}
	return sum.Sum, cursor.Decode(sum)
}

func (d *DatabaseBoard) Tick() error {
	ctx, cancel := newContext()
	defer cancel()
	if _, err := d.Collection.Aggregate(ctx, mongo.Pipeline{{{"$out", "copy"}}}); err != nil {
		return err
	}

	copiedCollection := d.Collection.Database().Collection("copy")
	for x := 0; x < BoardWidth; x++ {
		for y := 0; y < BoardHeight; y++ {
			sum, err := getSum(&d.Collection, x, y)
			if err != nil {
				return err
			}

			err = func() error {
				ctx, cancel := newContext()
				defer cancel()

				if sum < 2 || sum > 3 {
					_, err := copiedCollection.UpdateOne(
						ctx,
						bson.D{{"coordinate.0", x}, {"coordinate.1", y}},
						bson.M{"$set": bson.M{"cell": CellDead}},
					)
					if err != nil {
						return err
					}
				} else if sum == 3 {
					findResult := d.Collection.FindOne(
						ctx,
						bson.D{{"coordinate.0", x}, {"coordinate.1", y}},
					)
					if findResult.Err() != nil {
						return findResult.Err()
					}
					cellData := &CellData{}
					if err := findResult.Decode(cellData); err != nil {
						return err
					}

					if cellData.Cell == CellDead {
						_, err := copiedCollection.UpdateByID(
							ctx,
							cellData.ID,
							bson.M{"$set": bson.M{"cell": CellLive}},
						)
						if err != nil {
							return err
						}
					}
				}
				return nil
			}()
			if err != nil {
				return err
			}
		}
	}

	err := func() error {
		ctx, cancel := newContext()
		defer cancel()
		if err := d.Collection.Drop(ctx); err != nil {
			return err
		}

		_, err := copiedCollection.Aggregate(ctx, mongo.Pipeline{{{"$out", collectionName}}})
		return err
	}()
	return err
}
