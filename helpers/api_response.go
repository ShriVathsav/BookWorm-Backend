package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var ErrorMethodNotAllowed = "method Not allowed"

func ApiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{
		//"Content-Type":                 "application/json",
		"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "OPTIONS,POST,GET,DELETE,PUT",
	}}
	resp.StatusCode = status

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return ApiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
