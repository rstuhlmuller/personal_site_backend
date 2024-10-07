package utils

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// JSONResponse creates a new APIGatewayProxyResponse with JSON body
func JSONResponse(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	jsonBody, err := json.Marshal(body)
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

// ErrorResponse creates a new APIGatewayProxyResponse for error scenarios
func ErrorResponse(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	return JSONResponse(statusCode, map[string]string{"error": message})
}
