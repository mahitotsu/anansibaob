package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       "Hello Go lang",
	}, nil
}

func main() {
	lambda.Start(handler)
}
