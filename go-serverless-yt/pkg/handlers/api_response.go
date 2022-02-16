package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// http 요청을 처리해서 response를 보내주는 handler
func apiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status

	// go 언어에서 보통 정수형이나 구조체같은 타입을 바이트 슬라이스로 변경하는 것을 마셜링이라 한다.
	// go 언어는 json언어를 이해할 수 없기 때문에 이를 돕기위해 인코딩이 필요하다
	// 여기서는 body를 인코딩한 JSON값을 반환한다.
	// 인터페이스 값은 인터페이스에 포함된 값으로 인코딩된다.
	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}
