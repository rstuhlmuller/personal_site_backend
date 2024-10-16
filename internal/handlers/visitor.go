package handlers

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rstuhlmuller/personal_site_backend/internal/db"
	"github.com/rstuhlmuller/personal_site_backend/internal/models"
	"github.com/rstuhlmuller/personal_site_backend/pkg/utils"
)

func IncrementVisitorCount(ctx context.Context, req events.APIGatewayProxyRequest, dynamoDB db.DynamoDBInterface) (events.APIGatewayProxyResponse, error) {
	visitorInfo := models.NewVisitorLog(
		req.RequestContext.Identity.SourceIP,
		req.RequestContext.Identity.UserAgent,
		req.Headers["Referer"],
	)

	count, err := dynamoDB.IncrementVisitorCount(ctx, visitorInfo)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, "Failed to increment visitor count")
	}

	return utils.JSONResponse(http.StatusOK, map[string]int{"count": count})
}

func GetVisitorCount(ctx context.Context, req events.APIGatewayProxyRequest, dynamoDB db.DynamoDBInterface) (events.APIGatewayProxyResponse, error) {
	count, err := dynamoDB.GetVisitorCount(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, "Failed to get visitor count")
	}

	return utils.JSONResponse(http.StatusOK, map[string]int{"count": count})
}

func GetVisitorLog(ctx context.Context, req events.APIGatewayProxyRequest, dynamoDB db.DynamoDBInterface) (events.APIGatewayProxyResponse, error) {
	logs, err := dynamoDB.GetVisitorLog(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, "Failed to get visitor log")
	}

	return utils.JSONResponse(http.StatusOK, logs)
}
