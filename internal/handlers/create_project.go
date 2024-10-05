package handlers

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rstuhlmuller/personal_site_backend/internal/db"
	"github.com/rstuhlmuller/personal_site_backend/internal/models"
	"github.com/rstuhlmuller/personal_site_backend/pkg/utils"
)

func CreateProject(ctx context.Context, req events.APIGatewayProxyRequest, dynamoDB db.DynamoDBInterface) (events.APIGatewayProxyResponse, error) {
	var project models.Project
	err := json.Unmarshal([]byte(req.Body), &project)
	if err != nil {
		return utils.ErrorResponse(400, "Invalid request body")
	}

	err = dynamoDB.CreateProject(ctx, &project)
	if err != nil {
		return utils.ErrorResponse(500, "Failed to create project")
	}

	responseBody, _ := json.Marshal(project)
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(responseBody),
	}, nil
}
