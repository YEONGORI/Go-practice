package handlers

import (
	"go_src/custom-token/pkg/token"

	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
)

var ErrorMethodNotAllowed = "method not allowed"

type ErrorBody struct {
	// omitempty태그는 해당 필드(ErrorMsg)의 데이터가 존재하지 않으면, Marshaling할때 이 필드는 제외하고 Marshaling하게 한다.
	ErrorMsg *string `json:"error,omitempty"`
}

func GetToken(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	phoneNum := req.QueryStringParameters["phoneNum"]
	password := req.QueryStringParameters["password"]

	if len(phoneNum) > 0 || len(password) > 0 {
		result, err := token.GenerateToken(phoneNum, password)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
	} else {
		return apiResponse(http.StatusBadRequest, nil)
	}
}

func UnhandleMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, http.ErrBodyNotAllowed)
}
