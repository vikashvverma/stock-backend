package stock

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/vikashvverma/stock-backend/constants"
)

type Trader interface {
	Find(string) ([]PricePoint, error)
	FindAll([]string, time.Time, time.Time) ([]Stock, error)
	Top(time.Time, time.Time, bool) (interface{}, error)
}

type stockTrader struct {
	Client *mongo.Client
}

func New(c *mongo.Client) Trader {
	return &stockTrader{Client: c}
}

func (s *stockTrader) Find(name string) ([]PricePoint, error) {
	collection := s.Client.Database(constants.Database).Collection(constants.Collection)
	filter := bson.D{{
		Key: "$or",
		Value: bson.A{
			bson.D{{Key: "symbol", Value: bsonx.String(name)}},
			bson.D{{Key: "name", Value: bsonx.String(name)}},
		},
	}}
	res := collection.FindOne(context.Background(), filter)

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("find: error finding: %s", err)
	}

	var st Stock
	err := res.Decode(&st)
	if err != nil {
		return nil, fmt.Errorf("find: could not decode result: %s", err)
	}

	return st.PricePoints, nil
}

func (s *stockTrader) Top(from, to time.Time, best bool) (interface{}, error) {
	collection := s.Client.Database(constants.Database).Collection(constants.Collection)
	total := bson.D{
		{
			Key: "$sum",
			Value: bson.D{
				{Key: "$subtract", Value: bson.A{"$pricepoints.close", "$pricepoints.open"}},
			},
		},
	}

	order := 1
	if best {
		order = -1
	}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{}}},
		{{Key: "$unwind", Value: "$pricepoints"}},
		{{Key: "$match", Value: bson.D{{Key: "pricepoints.date", Value: bson.D{{Key: "$gte", Value: from}, {Key: "$lte", Value: to}}}}}},
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$symbol"}, {Key: "total", Value: total}}}},
		{{Key: "$sort", Value: bson.D{{Key: "total", Value: order}}}},
		{{Key: "$limit", Value: 10}},
	}

	ctx := context.Background()
	cur, err := collection.Aggregate(ctx, pipeline, options.Aggregate())
	if err != nil {
		return nil, fmt.Errorf("top: unable to find top stocks: %s", err)
	}
	defer cur.Close(ctx)

	var res []interface{}
	for cur.Next(ctx) {
		var result bson.M
		err = cur.Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("top: error decoding result: %s", err)
		}
		name, ok := result["_id"]
		if !ok {
			return nil, fmt.Errorf("top: unable to get name")
		}
		res = append(res, name)
	}

	return res, nil
}

func (s *stockTrader) FindAll(tickers []string, from, to time.Time) ([]Stock, error) {
	collection := s.Client.Database(constants.Database).Collection(constants.Collection)

	var names []interface{}
	for _, v := range tickers {
		names = append(names, v)
	}

	filter := bson.D{
		{
			Key: "symbol",
			Value: bson.D{
				{Key: "$in", Value: bson.A(names)},
			},
		}, {
			Key: "pricepoints.date", Value: bson.D{
				{Key: "$gte", Value: from},
				{Key: "$lte", Value: to},
			},
		},
	}

	ctx := context.Background()
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("findAll: unable to find stocks: %s", err)
	}
	defer cur.Close(ctx)

	var res []Stock
	for cur.Next(ctx) {
		var result Stock
		err = cur.Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("findAll: error decoding result: %s", err)
		}

		res = append(res, result)
	}

	return res, nil
}
