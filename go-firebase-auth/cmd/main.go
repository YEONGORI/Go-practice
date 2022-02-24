package main

import (
	"go_src/go-firebase-auth/pkg/handlers"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	lo, _ := zap.NewProduction()
	logger = lo
	defer logger.Sync()
}

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

	logger.Info("received event", zap.Any("method", req.HTTPMethod), zap.Any("path", req.Path), zap.Any("body", req.Body))

	var res *events.APIGatewayProxyResponse

	switch req.Path {
	case "/kakao/signin":
		return handlers.KakaoGetAuthCode(req)
	case "/kakao/token":
		return handlers.KakaoCreateCustomToken(req)
	case "/naver":
		return res, nil
	default:
		return res, nil
	}

}
