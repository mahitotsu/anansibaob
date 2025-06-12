package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v5"
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

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results, err := pgx.CollectRows(rows, pgx.RowToMap)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"results": results}, err
}

func main() {
	lambda.Start(handler)
}
