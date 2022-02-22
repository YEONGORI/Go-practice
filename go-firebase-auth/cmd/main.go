package main

import (
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
	// 임시방편으로 해둔 건데 각각 kakao,naver logic을 맡는 패키지로 분리해줘야함
	if req.Path == "/kakao" {
		return res, nil
	} else if req.Path == "/naver" {
		return res, nil
	} else {
		return res, nil
	}
	// switch req.HTTPMethod {
	// case "GET":
	// 	return handlers.GetAuthenticationCode(req)
	// case "POST":
	// 	return handlers.CreadteCustomtoken(req)
	// default:
	// 	return handlers.UnhandleMethod()
	// }
}
