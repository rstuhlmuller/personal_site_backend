package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/rstuhlmuller/personal_site_backend/internal/models"
)

type DynamoDBInterface interface {
	IncrementVisitorCount(ctx context.Context, visitorInfo *models.VisitorItem) (int, error)
	GetVisitorCount(ctx context.Context) (int, error)
	GetVisitorLog(ctx context.Context) ([]*models.VisitorItem, error) // Note the pointer here
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
	log.Printf("Attempting to increment visitor count and log visit in table: %s", db.table)

	// First, increment the count
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(db.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: models.CountItemID},
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

	updateResult, err := db.client.UpdateItem(ctx, updateInput)
	if err != nil {
		log.Printf("Error incrementing visitor count: %v", err)
		return 0, fmt.Errorf("failed to increment visitor count: %w", err)
	}

	var updatedCount models.VisitorItem
	err = attributevalue.UnmarshalMap(updateResult.Attributes, &updatedCount)
	if err != nil {
		log.Printf("Error unmarshalling updated count: %v", err)
		return 0, fmt.Errorf("failed to unmarshal updated count: %w", err)
	}

	// Now, log the visitor info
	visitorInfo.ID = uuid.New().String()
	visitorInfo.Timestamp = time.Now().UTC()

	item, err := attributevalue.MarshalMap(visitorInfo)
	if err != nil {
		log.Printf("Error marshalling visitor info: %v", err)
		return 0, fmt.Errorf("failed to marshal visitor info: %w", err)
	}

	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(db.table),
		Item:      item,
	}

	_, err = db.client.PutItem(ctx, putInput)
	if err != nil {
		log.Printf("Error logging visitor info: %v", err)
		return 0, fmt.Errorf("failed to log visitor info: %w", err)
	}

	log.Printf("Successfully incremented count to %d and logged visit", updatedCount.Count)

	return updatedCount.Count, nil
}

func (db *DynamoDB) GetVisitorCount(ctx context.Context) (int, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: models.CountItemID},
		},
	}

	result, err := db.client.GetItem(ctx, input)
	if err != nil {
		return 0, fmt.Errorf("failed to get visitor count: %v", err)
	}

	if result.Item == nil {
		// If the item doesn't exist, return 0 as the count
		return 0, nil
	}

	var count struct {
		Count int `dynamodbav:"count"`
	}
	err = attributevalue.UnmarshalMap(result.Item, &count)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal visitor count: %v", err)
	}

	return count.Count, nil
}

func (db *DynamoDB) GetVisitorLog(ctx context.Context) ([]*models.VisitorItem, error) {
	log.Printf("Attempting to get visitor log from table: %s", db.table)

	input := &dynamodb.ScanInput{
		TableName:        aws.String(db.table),
		FilterExpression: aws.String("#type = :logType"),
		ExpressionAttributeNames: map[string]string{
			"#type": "type",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":logType": &types.AttributeValueMemberS{Value: "log"},
		},
	}

	log.Printf("ScanInput: %+v", input)

	result, err := db.client.Scan(ctx, input)
	if err != nil {
		log.Printf("Error scanning table: %v", err)
		return nil, fmt.Errorf("failed to get visitor log: %w", err)
	}

	log.Printf("Scan result: %+v", result)

	if len(result.Items) == 0 {
		log.Printf("No items found in scan result")
		return []*models.VisitorItem{}, nil
	}

	var logs []*models.VisitorItem
	err = attributevalue.UnmarshalListOfMaps(result.Items, &logs)
	if err != nil {
		log.Printf("Error unmarshalling visitor logs: %v", err)
		return nil, fmt.Errorf("failed to unmarshal visitor logs: %w", err)
	}

	log.Printf("Successfully retrieved %d log entries", len(logs))

	return logs, nil
}
