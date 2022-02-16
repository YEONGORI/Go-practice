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
	ErrorCouldNotMarshalItem     = "could not marshal item"
	ErrorCouldNotDeleteItem      = "could not delete item"
	ErrorCouldNotDynamoPutItem   = "could not dynamo put item"
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

	// 여기서 User객체의 요소(eleement) 타입이 String에서 AttributeValue로 변한다.
	// 115라인에서 PutItemInput의 required field인 Item을 초기화 하기 위해
	av, err := dynamodbattribute.MarshalMap(u)

	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &u, nil
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var u User

	// Unmarshal이란 byte형식에서 다른 형식으로 바뀌는 것이다.
	// marshaling의 방향은 byte로의 단방향 변환이지만 unmarshaling은 단방향이 아니다. 따라서 유저가 직접 방향을 지정해줄 필요가 있다.
	// 여기서는 req.Body에 담긴 []byte타입('string타입으로 생각할 수 있음')의 데이터를 u(User) 타입으로 변환해준다.
	// 또 두번째 인자는 반드시 포인터 여야한다. 그렇지 않다면 단순한 값 복사만 일어나 변환되지 않는다.
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidEmail)
	}

	currentUser, _ := FetchUser(u.Email, tableName, dynaClient)
	if currentUser != nil && len(currentUser.Email) == 0 {
		return nil, errors.New(ErrorUserDoesNotExist)
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &u, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		// 삭제할 item의 기본 키를 나타내는 attribute의 이름(email)을 AttributeValue 개체에 매핑한다.
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	return nil
}
