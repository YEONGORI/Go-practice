package handlers

import (
	"custom-token/pkg/token"

	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
)

var c = "qwe"
var ErrorMethodNotAllowed = "method not allowed"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitEmpty"`
}

func GetToken(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	phoneNum := req.QueryStringParameters["phoneNum"]
	password := req.QueryStringParameters["password"]

	if len(phoneNum) > 0 || len(password) > 0 {
		result, err := token.FetchToken(phoneNum, password)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
	} else {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
}
