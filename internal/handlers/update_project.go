package handlers

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rstuhlmuller/personal_site_backend/internal/db"
	"github.com/rstuhlmuller/personal_site_backend/pkg/utils"
)

func UpdateProject(ctx context.Context, req events.APIGatewayProxyRequest, dynamoDB db.DynamoDBInterface) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]
	if id == "" {
		return utils.ErrorResponse(400, "Missing project ID")
	}

	var updates map[string]interface{}
	err := json.Unmarshal([]byte(req.Body), &updates)
	if err != nil {
		return utils.ErrorResponse(400, "Invalid request body")
	}

	project, err := dynamoDB.GetProject(ctx, id)
	if err != nil {
		return utils.ErrorResponse(500, "Failed to get project")
	}
	if project == nil {
		return utils.ErrorResponse(404, "Project not found")
	}

	project.Update(updates)

	err = dynamoDB.UpdateProject(ctx, project)
	if err != nil {
		return utils.ErrorResponse(500, "Failed to update project")
	}

	return utils.JSONResponse(200, project)
}
