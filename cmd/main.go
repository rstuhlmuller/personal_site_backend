package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rstuhlmuller/personal_site_backend/internal/db"
	"github.com/rstuhlmuller/personal_site_backend/internal/handlers"
)

var dynamoDB db.DynamoDBInterface

func init() {
	var err error
	dynamoDB, err = db.NewDynamoDB()
	if err != nil {
		log.Fatalf("Failed to initialize DynamoDB: %v", err)
	}
}

func addCorsHeaders(resp *events.APIGatewayProxyResponse) {
	if resp.Headers == nil {
		resp.Headers = make(map[string]string)
	}
	resp.Headers["Access-Control-Allow-Origin"] = "https://rodman.stuhlmuller.net"
	resp.Headers["Access-Control-Allow-Methods"] = "GET,POST,OPTIONS"
	resp.Headers["Access-Control-Allow-Headers"] = "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token"
}

func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var resp events.APIGatewayProxyResponse

	// Handle CORS preflight requests
	if req.HTTPMethod == "OPTIONS" {
		resp = events.APIGatewayProxyResponse{StatusCode: 200}
		addCorsHeaders(&resp)
		return resp, nil
	}

	switch req.HTTPMethod {
	case "GET":
		if req.Resource == "/visitor_count" {
			countResp, err := handlers.GetVisitorCount(ctx, req, dynamoDB)
			if err != nil {
				log.Printf("Error getting visitor count: %v", err)
				resp = events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}
			} else {
				var visitorCount map[string]int
				err := json.Unmarshal([]byte(countResp.Body), &visitorCount)
				if err != nil {
					log.Printf("Error parsing visitor count: %v", err)
				}
				count := visitorCount["count"]
				body, _ := json.Marshal(map[string]int{"count": count})
				resp = events.APIGatewayProxyResponse{StatusCode: 200, Body: string(body)}
			}
		}
		// ... other GET routes ...
	case "POST":
		if req.Resource == "/visitor_count" {
			countResp, err := handlers.IncrementVisitorCount(ctx, req, dynamoDB)
			if err != nil {
				log.Printf("Error incrementing visitor count: %v", err)
				resp = events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}
			} else {
				var visitorCount map[string]int
				err := json.Unmarshal([]byte(countResp.Body), &visitorCount)
				if err != nil {
					log.Printf("Error parsing visitor count: %v", err)
				}
				count := visitorCount["count"]
				body, _ := json.Marshal(map[string]int{"count": count})
				resp = events.APIGatewayProxyResponse{StatusCode: 200, Body: string(body)}
			}
		}
	default:
		resp = events.APIGatewayProxyResponse{StatusCode: 404, Body: "Not Found"}
	}

	addCorsHeaders(&resp)
	return resp, nil
}

func main() {
	lambda.Start(router)
}
