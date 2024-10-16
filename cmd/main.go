package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rstuhlmuller/personal_site_backend/internal/db"
	"github.com/rstuhlmuller/personal_site_backend/internal/handlers"
)

var dynamoDB db.DynamoDBInterface
var allowedOrigins map[string]bool

func init() {
	log.Printf("Initializing")
	var err error
	dynamoDB, err = db.NewDynamoDB()
	if err != nil {
		log.Fatalf("Failed to initialize DynamoDB: %v", err)
	}

	allowedOrigins = make(map[string]bool)
	origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	for _, origin := range origins {
		allowedOrigins[strings.TrimSpace(origin)] = true
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

	log.Printf("Request: %v", req)
	// Handle CORS preflight requests
	if req.HTTPMethod == "OPTIONS" {
		resp = events.APIGatewayProxyResponse{StatusCode: 200}
		addCorsHeaders(&resp)
		return resp, nil
	}

	log.Printf("Request HTTP Method: %v", req.HTTPMethod)
	switch req.HTTPMethod {
	case "GET":
		log.Printf("Request resource: %v", req.Path)
		if req.Path == "/visitor_count" {
			countResp, err := handlers.GetVisitorCount(ctx, req, dynamoDB)
			if err != nil {
				log.Printf("Error getting visitor count: %v", err)
				resp = events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}
			} else {
				var visitorCount map[string]int
				log.Printf("Unmarshalling visitor count: %v", countResp.Body)
				err := json.Unmarshal([]byte(countResp.Body), &visitorCount)
				if err != nil {
					log.Printf("Error parsing visitor count: %v", err)
				}
				count := visitorCount["count"]
				body, _ := json.Marshal(map[string]int{"count": count})
				resp = events.APIGatewayProxyResponse{StatusCode: 200, Body: string(body)}
			}
		}
		if req.Path == "/visitor_log" {
			return handlers.GetVisitorLog(ctx, req, dynamoDB)
		}
		// ... other GET routes ...
	case "POST":
		log.Printf("Request resource: %v", req.Path)
		if req.Path == "/visitor_count" {
			log.Printf("Incrementing visitor count: %v", req.RequestContext.Identity.SourceIP)
			countResp, err := handlers.IncrementVisitorCount(ctx, req, dynamoDB)
			if err != nil {
				log.Printf("Error incrementing visitor count: %v", err)
				resp = events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}
			} else {
				var visitorCount map[string]int
				log.Printf("Unmarshalling visitor count: %v", countResp.Body)
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
