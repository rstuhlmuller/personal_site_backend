package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rstuhlmuller/personal_site_backend/internal/models"
)

type DynamoDBInterface interface {
	IncrementVisitorCount(ctx context.Context, visitorInfo *models.VisitorItem) (int, error)
	GetVisitorCount(ctx context.Context) (int, error)
	GetVisitorLog(ctx context.Context) ([]*models.VisitorItem, error)
}

type DynamoDB struct {
	client *dynamodb.Client
	table  string
}

func NewDynamoDB() (*DynamoDB, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}

	tableName := os.Getenv("DYNAMODB_TABLE")
	log.Printf("DYNAMODB_TABLE: %v", tableName)
	if tableName == "" {
		return nil, fmt.Errorf("DYNAMODB_TABLE environment variable is not set")
	}

	return &DynamoDB{
		client: dynamodb.NewFromConfig(cfg),
		table:  tableName,
	}, nil
}

func (db *DynamoDB) IncrementVisitorCount(ctx context.Context, visitorInfo *models.VisitorItem) (int, error) {
	log.Printf("Attempting to increment visitor count in table: %s", db.table)

	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(db.table),
		Key: map[string]types.AttributeValue{
			"id":   &types.AttributeValueMemberS{Value: models.CountItemID},
			"type": &types.AttributeValueMemberS{Value: "count"},
		},
		UpdateExpression: aws.String("SET #count = if_not_exists(#count, :zero) + :incr, #timestamp = :now"),
		ExpressionAttributeNames: map[string]string{
			"#count":     "count",
			"#timestamp": "timestamp",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":incr": &types.AttributeValueMemberN{Value: "1"},
			":zero": &types.AttributeValueMemberN{Value: "0"},
			":now":  &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	}

	log.Printf("UpdateInput: %+v", updateInput)
	log.Printf("UpdateInput Key: %+v", updateInput.Key)
	log.Printf("UpdateInput ExpressionAttributeNames: %+v", updateInput.ExpressionAttributeNames)
	log.Printf("UpdateInput ExpressionAttributeValues: %+v", updateInput.ExpressionAttributeValues)

	updateResult, err := db.client.UpdateItem(ctx, updateInput)
	if err != nil {
		log.Printf("Error incrementing visitor count: %v", err)
		// Check for specific error types
		var notFoundErr *types.ResourceNotFoundException
		var conditionalCheckFailedErr *types.ConditionalCheckFailedException
		if errors.As(err, &notFoundErr) {
			log.Printf("Table or item not found: %v", notFoundErr)
		} else if errors.As(err, &conditionalCheckFailedErr) {
			log.Printf("Conditional check failed: %v", conditionalCheckFailedErr)
		}
		return 0, fmt.Errorf("failed to increment visitor count: %w", err)
	}

	log.Printf("Update result: %+v", updateResult)

	var updatedCount models.VisitorItem
	err = attributevalue.UnmarshalMap(updateResult.Attributes, &updatedCount)
	if err != nil {
		log.Printf("Error unmarshalling updated count: %v", err)
		return 0, fmt.Errorf("failed to unmarshal updated count: %w", err)
	}

	log.Printf("Updated count: %d", updatedCount.Count)

	return updatedCount.Count, nil
}

func (db *DynamoDB) GetVisitorCount(ctx context.Context) (int, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.table),
		Key: map[string]types.AttributeValue{
			"id":   &types.AttributeValueMemberS{Value: models.CountItemID},
			"type": &types.AttributeValueMemberS{Value: "count"},
		},
	}

	result, err := db.client.GetItem(ctx, input)
	if err != nil {
		return 0, fmt.Errorf("failed to get visitor count: %v", err)
	}

	if result.Item == nil {
		return 0, nil
	}

	var count models.VisitorItem
	err = attributevalue.UnmarshalMap(result.Item, &count)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal visitor count: %v", err)
	}

	return count.Count, nil
}

func (db *DynamoDB) GetVisitorLog(ctx context.Context) ([]*models.VisitorItem, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.table),
		KeyConditionExpression: aws.String("#type = :logType"),
		ExpressionAttributeNames: map[string]string{
			"#type": "type",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":logType": &types.AttributeValueMemberS{Value: "log"},
		},
	}

	result, err := db.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get visitor log: %v", err)
	}

	var logs []*models.VisitorItem
	err = attributevalue.UnmarshalListOfMaps(result.Items, &logs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal visitor logs: %v", err)
	}

	return logs, nil
}
