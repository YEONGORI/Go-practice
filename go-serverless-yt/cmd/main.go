/*
Lambda란 서버를 신경쓰지 않고도 코드를 실행할 수 있게 해주는 서비스이다.
Lambda서비스에 작성한 코드를 올려두고 HTTP요청이나 AWS내의 이벤트가 발생하면
코드가 실행되는 구조이다.

간단하게 말하자면 서버리스로 작동하는 API관리/실행 서비스이다.
*/
package main

import (
	"go-serverless-yt/pkg/handlers"
	"os"

	// 핸들러가 처리할 이벤트 리소스들의 유형을 정의한 패키지
	"github.com/aws/aws-lambda-go/events"
	// AWS Lambda에서 핸들러를 호출하기 위해 사용되는 패키지이자
	// Go가 Lambda에서 실행되도록하는 패키지
	"github.com/aws/aws-lambda-go/lambda"
	// go 언어로 aws 서비스를 사용하기 위한 SDK(소프트웨어 개발 키트)
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	// Amazon DynamoDB서비스 클라이언트를 모킹하여 코드를 테스트할 수 있는
	// 인터페이스를 제공한다.  모킹이란? -> 테스트하려는 코드에 외부 의존성이
	// 있을 때 unit테스트에서 사용되는 프로세스다. 실제 객체를 대체 객체로 바꾼다.
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynaClient dynamodbiface.DynamoDBAPI
)

// aws-sdk-go 는 Go application을 Amazon S3, Amazon DynamoDB등을 비롯한
// AWS 서비스 제품들과 통합하기 위해 사용한다.
func main() {
	region := os.Getenv("AWS_REGION")
	// 모든 aws client는 Session을 필요로 한다. Session은 client에게
	// 사용지역, endpoint, credential(자격증명)같은 Configuration을 제공한다.
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return
	}

	// DynamoDB의 client(DynamoDB의 인스턴스)를 생성한다.
	dynaClient = dynamodb.New(awsSession)

	// 이 어플리케이션은 server가 존재하지 않고 lambda로 동작하기 때문에
	// lambda.Start가 반드시 필요하다.
	lambda.Start(handler)
}

const tableName = "LambdaInGoUser"

// Lambda 함수의 핸들러는 이벤트(events)를 처리하는 함수 코드의 메서드이다.
// APIGateway로 부터 HTTP요청을 전달 받아서 무언가를 처리하고 HTTP응답의 형태로 돌려준다.
func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return handlers.GetUser(req, tableName, dynaClient)
	case "POST":
		return handlers.CreateUser(req, tableName, dynaClient)
	case "PUT":
		return handlers.UpdateUser(req, tableName, dynaClient)
	case "DELETE":
		return handlers.DeleteUser(req, tableName, dynaClient)
	default:
		return handlers.UnhandleMethod()
	}

}
