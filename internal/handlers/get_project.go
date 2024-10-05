package handlers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rstuhlmuller/personal_site_backend/internal/db"
	"github.com/rstuhlmuller/personal_site_backend/pkg/utils"
)

func GetProject(ctx context.Context, req events.APIGatewayProxyRequest, dynamoDB db.DynamoDBInterface) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]
	if id == "" {
		return utils.ErrorResponse(400, "Missing project ID")
	}

	project, err := dynamoDB.GetProject(ctx, id)
	if err != nil {
		return utils.ErrorResponse(500, "Failed to get project")
	}
	if project == nil {
		return utils.ErrorResponse(404, "Project not found")
	}

	return utils.JSONResponse(200, project)
}
