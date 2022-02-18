package main

import (
	"go_src/custom-token/pkg/handlers"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {
	region := os.Getenv("AWS_REGION")

	_, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return
	}

	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return handlers.GetToken(req)
	default:
		return handlers.UnhandleMethod()
	}
}

// func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyRequest, error) {
// 	switch req.HTTPMethod {
// 	case "GET":
// 		return req, nil
// 	default:
// 		return req, nil
// 	}
// }
