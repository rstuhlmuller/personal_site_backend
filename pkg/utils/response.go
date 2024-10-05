package utils

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type APIResponse struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

func JSONResponse(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	response := APIResponse{
		StatusCode: statusCode,
		Body:       body,
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(jsonBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func ErrorResponse(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	return JSONResponse(statusCode, map[string]string{"error": message})
}
