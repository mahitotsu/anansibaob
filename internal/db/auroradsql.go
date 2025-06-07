package db

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
)

func GenerateDbConnectAdminAuthToken(clusterEndpoint string, region string, action string) (string, error) {

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("config loading failed: %w", err)
	}

	endpoint := "https://" + clusterEndpoint
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("request creation failed: %w", err)
	}

	q := req.URL.Query()
	q.Add("Action", action)
	req.URL.RawQuery = q.Encode()

	signer := v4.NewSigner()
	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return "", fmt.Errorf("credentials retrieval failed: %w", err)
	}

	signedURL, _, err := signer.PresignHTTP(
		ctx,
		creds,
		req,
		"",
		"dsql",
		region,
		time.Now().Add(15*time.Minute),
	)
	if err != nil {
		return "", fmt.Errorf("signing failed: %w", err)
	}

	return strings.TrimPrefix(signedURL, "https://"), nil
}
