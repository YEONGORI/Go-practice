package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func apiResponse(base64Encodeed bool, status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status
	resp.IsBase64Encoded = false
	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}
