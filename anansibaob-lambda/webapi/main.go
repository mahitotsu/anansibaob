package main

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello go lang function.",
		})
	})

	ginLambda = ginadapter.New(r)
}

func handler(ctx context.Context, funcReq events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	proxyReq := events.APIGatewayProxyRequest{
		HTTPMethod:            funcReq.RequestContext.HTTP.Method,
		Path:                  funcReq.RequestContext.HTTP.Path,
		Headers:               funcReq.Headers,
		Body:                  funcReq.Body,
		IsBase64Encoded:       funcReq.IsBase64Encoded,
		QueryStringParameters: funcReq.QueryStringParameters,
	}

	proxyRes, err := ginLambda.ProxyWithContext(ctx, proxyReq)

	headers := make(map[string]string)
	for key, values := range proxyRes.MultiValueHeaders {
		headers[key] = strings.Join(values, ",")
	}
	for key, value := range proxyRes.Headers {
		headers[key] = value
	}

	funcRes := events.LambdaFunctionURLResponse{
		StatusCode:      proxyRes.StatusCode,
		Headers:         headers,
		Body:            proxyRes.Body,
		IsBase64Encoded: proxyRes.IsBase64Encoded,
	}

	return funcRes, err
}

func main() {
	lambda.Start(handler)
}
