package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v5"
	"github.com/mahitotsu/anansibaob/internal/db"
)

func handler(ctx context.Context, event map[string]interface{}) (map[string]interface{}, error) {

	context := context.Background()
	endpoint := os.Getenv("DB_CLUSTER_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("DB_CLUSTER_ENDPOINT environment variable not set")
	}

	parts := strings.Split(endpoint, ".")
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid endpoint format: %s", endpoint)
	}
	region := parts[2]

	token, err := db.GenerateDbConnectAdminAuthToken(endpoint, region, "DbConnectAdmin")
	if err != nil {
		return nil, err
	}

	dbUrl := fmt.Sprintf("postgres://%s:5432/postgres?user=admin&sslmode=require", endpoint)
	dbConfig, err := pgx.ParseConfig(dbUrl)
	dbConfig.Password = token

	conn, err := pgx.ConnectConfig(context, dbConfig)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context)

	rows, err := conn.Query(ctx, "SELECT 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		rowMap, err := pgx.RowToMap(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, rowMap)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return map[string]interface{}{"results": results}, nil
}

func main() {
	lambda.Start(handler)
}
