package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mahitotsu/anansibaob/internal/dao"
)

var sqlClient dao.SqlClient

func init() {

	clusterEndpoint := os.Getenv("DB_CLUSTER_ENDPOINT")
	if clusterEndpoint == "" {
		panic(fmt.Errorf("DB_CLUSTER_ENDPOINT environment variable not set"))
	}

	pool, err := dao.CreatePgPool(context.Background(), clusterEndpoint)
	if err != nil {
		panic(fmt.Errorf("failed to create database connection pool: %w", err))
	}

	sqlClient = dao.NewSqlClient(pool)
}

func handler(ctx context.Context, event map[string]interface{}) (map[string]interface{}, error) {

	results, err := sqlClient.Query(ctx, "SELECT 1")
	return map[string]interface{}{"results": results}, err
}

func main() {
	lambda.Start(handler)
}
