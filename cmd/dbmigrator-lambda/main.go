package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mahitotsu/anansibaob/internal/dao"
)

func handler(ctx context.Context, event map[string]interface{}) (map[string]interface{}, error) {

	endpoint := os.Getenv("DB_CLUSTER_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("DB_CLUSTER_ENDPOINT environment variable not set")
	}

	db, err := dao.OpenPgxDb(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"results": db.Config().Host}, err
}

func main() {
	lambda.Start(handler)
}
