/*
이곳에 있는 함수들은 handlers.go 파일의 함수들과 1:1 대응한다.
main함수의 handler(router가 더 가까운 표현일듯) 에서 각 요청이 들어오면
handlers.go파일의 함수들을 호출하고 그 함수들이 이곳(user.go)파일의 함수를
호출하면 이곳의 함수들이 Database과 직접 관계해 응답을 주고받는다.
*/
package user

import (
	"encoding/json"
	"errors"
	"go-serverless-yt/pkg/validators"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrorFailedToFetchRecord     = "failed to fetch record"
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorInvalidUserData         = "invalid user data"
	ErrorInvalidEmail            = "invalid email"
	ErrorCouldNotMarshal         = "could not marshal item"
	ErrorCouldNotDeleteItem      = "could not delete item"
	ErrorCouldNotDynamoPuItem    = "could not dynamo put item"
	ErrorUserAlreadyExists       = "user.User already exists"
	ErrorUserDoesNotExist        = "user.User does not exist"
)

// Golang에서는 json 직렬화(serialize, marshal)할 때 구조체 내 필드들의 첫
// 글자를 대문자로 쓰지 않으면 직렬화할 수 없다. 이때 직렬화된 json의 Key값을
// 소문자로 하기 위해서는 Tag(`` 백틱)을 사용해 key값을 지정해야한다.
type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func FetchUser(email, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	// GetItemInput은 GetItem 작업의 입력을 나타내는 메소드
	// 그러므로 변수input은 query라 볼 수 있다.
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email), // S는 String type의 attribute
			},
		},
		// TableName은 요청된 item을 담고있는 table의 이름
		TableName: aws.String(tableName),
	}

	// DynamoDB로 부터 data(Item)를 가져온다.
	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	// dynamodb.AttributeValue 타입을 go 타입으로 unmarshal한다.
	// 첫번째 인자로 전자의 타입을 가진 Item을, 두번쨰 인자로 값을 담을 go타입을 넣어준다.
	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)

	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}

func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]User, error) {
	// ScanInput은 구조체다.ScanInput 구조체 내의 TableName 필드는 requried field이므로 반드시 포함해야한다.
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	// Scan메소드는 주어진 Table이나 인덱스의 모든 항목에 접근하여 그것에 해당하는 모든 item을 반환한다.
	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)

	return item, nil
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var u User

	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	if !validators.IsEmailValid(u.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}

	currentUser, _ := FetchUser(u.Email, tableName, dynaClient)
	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExists)
	}
	// 여기까지 1시간 8분@@@@@@@@@@@@@@@@@@
	// 여기까지 1시간 8분@@@@@@@@@@@@@@@@@@
	// 여기까지 1시간 8분@@@@@@@@@@@@@@@@@@
	// 여기까지 1시간 8분@@@@@@@@@@@@@@@@@@
	// 여기까지 1시간 8분@@@@@@@@@@@@@@@@@@
	// 여기까지 1시간 8분@@@@@@@@@@@@@@@@@@
	// 여기까지 1시간 8분@@@@@@@@@@@@@@@@@@

}

func UpdateUser() {}

func DeleteUser() error {}
