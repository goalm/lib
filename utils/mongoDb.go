package utils

import (
	"context"

	"github.com/qiniu/qmgo"
)

// GetMongoClient returns a new mongo client
// Candidate for storing model results
// How to use:
//		cli := utils.GetMongoClient()
//		defer func() {
//			if err := cli.Close(ctx); err != nil {
//				panic(err)
//			}
//		}()
//	 result, err := cli.Database("test").Collection("user").InsertOne(ctx, userInfo)

func GetMongoClient(ctx context.Context) *qmgo.Client {
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017"})
	if err != nil {
		panic(err)
	}
	return client
}
