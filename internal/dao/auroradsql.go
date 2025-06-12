package dao

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dsql/auth"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePgPool(ctx context.Context, clusterEndpoint string) (*pgxpool.Pool, error) {

	parts := strings.Split(clusterEndpoint, ".")
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid endpoint format: %s", clusterEndpoint)
	}
	region := parts[2]

	token, err := GenerateDbConnectAdminAuthToken(ctx, clusterEndpoint, region)
	if err != nil {
		return nil, err
	}

	dbUrl := fmt.Sprintf("postgres://%s:5432/postgres?user=admin&sslmode=require", clusterEndpoint)
	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, err
	}
	dbConfig.ConnConfig.Password = token

	return pgxpool.NewWithConfig(ctx, dbConfig)
}

func GenerateDbConnectAdminAuthToken(ctx context.Context, clusterEndpoint, region string) (string, error) {

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}

	return auth.GenerateDBConnectAdminAuthToken(ctx, clusterEndpoint, region, cfg.Credentials)
}
