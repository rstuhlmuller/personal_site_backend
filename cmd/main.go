package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rstuhlmuller/personal_site_backend/internal/db"
	"github.com/rstuhlmuller/personal_site_backend/internal/handlers"
	"github.com/rstuhlmuller/personal_site_backend/pkg/utils"
)

func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamoDB, err := db.NewDynamoDB()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	switch req.HTTPMethod {
	case "GET":
		if req.Resource == "/visitor-count" {
			return handlers.GetVisitorCount(ctx, req, dynamoDB)
		}
		// ... other GET routes ...
	case "POST":
		if req.Resource == "/visitor-count" {
			return handlers.IncrementVisitorCount(ctx, req, dynamoDB)
		}
		// ... other POST routes ...
	}

	return utils.ErrorResponse(404, "Not Found")
}

func main() {
	lambda.Start(router)
}
