package db

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

func GenerateDbConnectAdminAuthToken(clusterEndpoint string, region string, action string) (string, error) {

	sess, err := session.NewSession()
	if err != nil {
		return "", fmt.Errorf("session creation failed: %w", err)
	}

	endpoint := "https://" + clusterEndpoint + "/"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("request creation failed: %w", err)
	}

	q := req.URL.Query()
	q.Add("Action", action)
	req.URL.RawQuery = q.Encode()

	signer := v4.NewSigner(sess.Config.Credentials)
	_, err = signer.Presign(req, nil, "dsql", region, 15*time.Minute, time.Now())
	if err != nil {
		return "", fmt.Errorf("signing failed: %w", err)
	}

	return strings.TrimPrefix(req.URL.String(), "https://"), nil
}
