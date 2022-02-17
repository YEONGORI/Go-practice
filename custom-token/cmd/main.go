package main

import (
	//"go_src/custom-token/pkg/handlers"
	"os"

	// 핸들러가 처리할 이벤트 리소스들의 유형을 정의한 패키지
	"github.com/aws/aws-lambda-go/events"

	// AWS Lambda에서 핸들러를 호출하기 위해 사용되는 패키지이자
	// Go가 Lambda에서 실행되도록하는 패키
	"github.com/aws/aws-lambda-go/lambda"

	// go 언어로 aws 서비스를 사용하기 위한 SDK(소프트웨어 개발 키트)
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

// func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
// 	switch req.HTTPMethod {
// 	case "GET":
// 		return handlers.GetToken(req)
// 	default:
// 		return handlers.UnhandleMethod()
// 	}
// }

func handler(req events.APIGatewayProxyRequest) (string, error) {
	switch req.HTTPMethod {
	case "GET":
		return "hi", nil
	default:
		return req.HTTPMethod, nil
	}
}
