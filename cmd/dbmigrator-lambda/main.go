package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mahitotsu/anansibaob/internal/db"
)

func handler(ctx context.Context, event map[string]interface{}) (map[string]interface{}, error) {

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

	return map[string]interface{}{"token": token}, nil
}

func main() {
	lambda.Start(handler)
}
