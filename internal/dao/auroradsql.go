package dao

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

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

	poolConfig, err := pgxpool.ParseConfig("postgres://")
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgxpool config: %w", err)
	}
	poolConfig.ConnConfig.Host = clusterEndpoint
	poolConfig.ConnConfig.Port = 5432
	poolConfig.ConnConfig.User = "admin"
	poolConfig.ConnConfig.Password = token
	poolConfig.ConnConfig.Database = "postgres"
	poolConfig.ConnConfig.TLSConfig = &tls.Config{
		ServerName:         clusterEndpoint,
		InsecureSkipVerify: true,
	}

	poolConfig.MaxConnIdleTime = 5 * time.Minute
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.MaxConns = 1
	poolConfig.MinConns = 0

	return pgxpool.NewWithConfig(ctx, poolConfig)
}

func GenerateDbConnectAdminAuthToken(ctx context.Context, clusterEndpoint, region string) (string, error) {

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}

	return auth.GenerateDBConnectAdminAuthToken(ctx, clusterEndpoint, region, cfg.Credentials)
}
