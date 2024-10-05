package handlers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rstuhlmuller/personal_site_backend/internal/db"
	"github.com/rstuhlmuller/personal_site_backend/pkg/utils"
)

func ListProjects(ctx context.Context, req events.APIGatewayProxyRequest, dynamoDB db.DynamoDBInterface) (events.APIGatewayProxyResponse, error) {
	projects, err := dynamoDB.ListProjects(ctx)
	if err != nil {
		return utils.ErrorResponse(500, "Failed to list projects")
	}

	return utils.JSONResponse(200, projects)
}
