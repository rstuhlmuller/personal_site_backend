package handlers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rstuhlmuller/personal_site_backend/internal/db"
	"github.com/rstuhlmuller/personal_site_backend/pkg/utils"
)

func IncrementVisitorCount(ctx context.Context, req events.APIGatewayProxyRequest, dynamoDB db.DynamoDBInterface) (events.APIGatewayProxyResponse, error) {
	count, err := dynamoDB.IncrementVisitorCount(ctx)
	if err != nil {
		return utils.ErrorResponse(500, "Failed to increment visitor count")
	}

	return utils.JSONResponse(200, map[string]int{"count": count})
}

func GetVisitorCount(ctx context.Context, req events.APIGatewayProxyRequest, dynamoDB db.DynamoDBInterface) (events.APIGatewayProxyResponse, error) {
	count, err := dynamoDB.GetVisitorCount(ctx)
	if err != nil {
		return utils.ErrorResponse(500, "Failed to get visitor count")
	}

	return utils.JSONResponse(200, map[string]int{"count": count})
}
