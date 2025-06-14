package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mahitotsu/anansibaob/internal/dao"
)

var dbInitializer dao.DBInitializer

func init() {

	clusterEndpoint := os.Getenv("DB_CLUSTER_ENDPOINT")
	if clusterEndpoint == "" {
		panic(fmt.Errorf("DB_CLUSTER_ENDPOINT environment variable not set"))
	}

	if pool, err := dao.CreatePgPool(context.Background(), clusterEndpoint); err != nil {
		panic(fmt.Errorf("failed to create database connection pool: %w", err))
	} else {
		dbInitializer = dao.NewDBInitializer(dao.NewSqlClient(pool))
	}
}

func handler(ctx context.Context, event map[string]interface{}) (map[string]interface{}, error) {

	if err := dbInitializer.DropAndCreate(ctx); err != nil {
		return nil, err
	}
	return map[string]interface{}{"status": "success"}, nil
}

func main() {
	lambda.Start(handler)
}
