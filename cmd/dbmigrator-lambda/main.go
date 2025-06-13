package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahitotsu/anansibaob/internal/dao"
)

var pool *pgxpool.Pool
var initerr error

func init() {

	clusterEndpoint := os.Getenv("DB_CLUSTER_ENDPOINT")
	if clusterEndpoint == "" {
		panic(fmt.Errorf("DB_CLUSTER_ENDPOINT environment variable not set"))
	}

	pool, initerr = dao.CreatePgPool(context.Background(), clusterEndpoint)
}

func handler(ctx context.Context, event map[string]interface{}) (map[string]interface{}, error) {

	if initerr != nil {
		return nil, initerr
	}

	results, err := dao.NewSqlClient(pool).Query(ctx, "SELECT 1")

	return map[string]interface{}{"results": results}, err
}

func main() {
	lambda.Start(handler)
}
