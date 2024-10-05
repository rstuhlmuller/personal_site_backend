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
	dynamoDB := db.NewDynamoDB()

	switch req.HTTPMethod {
	case "GET":
		if req.Resource == "/projects" {
			return handlers.ListProjects(ctx, req, dynamoDB)
		} else if req.Resource == "/projects/{id}" {
			return handlers.GetProject(ctx, req, dynamoDB)
		}
	case "POST":
		if req.Resource == "/projects" {
			return handlers.CreateProject(ctx, req, dynamoDB)
		}
	case "PUT":
		if req.Resource == "/projects/{id}" {
			return handlers.UpdateProject(ctx, req, dynamoDB)
		}
	}

	return utils.ErrorResponse(404, "Not Found")
}

func main() {
	lambda.Start(router)
}
