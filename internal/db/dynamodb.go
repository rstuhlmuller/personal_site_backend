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
	"github.com/rstuhlmuller/personal_site_backend/internal/models"
)

type DynamoDBInterface interface {
	IncrementVisitorCount(ctx context.Context) (int, error)
	GetVisitorCount(ctx context.Context) (int, error)
	// Add other method signatures as needed
}

type DynamoDB struct {
	client *dynamodb.Client
	table  string
}

func NewDynamoDB() (*DynamoDB, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	tableName := os.Getenv("DYNAMODB_TABLE")
	if tableName == "" {
		return nil, fmt.Errorf("DYNAMODB_TABLE environment variable is not set")
	}

	return &DynamoDB{
		client: client,
		table:  tableName,
	}, nil
}

func (db *DynamoDB) IncrementVisitorCount(ctx context.Context) (int, error) {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(db.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: "visitor_count"},
		},
		UpdateExpression: aws.String("SET #count = if_not_exists(#count, :zero) + :incr, #date = :now"),
		ExpressionAttributeNames: map[string]string{
			"#count": "count",
			"#date":  "date",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":incr": &types.AttributeValueMemberN{Value: "1"},
			":zero": &types.AttributeValueMemberN{Value: "0"},
			":now":  &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	}

	result, err := db.client.UpdateItem(ctx, input)
	if err != nil {
		return 0, fmt.Errorf("failed to increment visitor count: %v", err)
	}

	var updatedCount models.VisitorCount
	err = attributevalue.UnmarshalMap(result.Attributes, &updatedCount)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal updated count: %v", err)
	}

	return updatedCount.Count, nil
}

func (db *DynamoDB) GetVisitorCount(ctx context.Context) (int, error) {
	result, err := db.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(db.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: "visitor_count"},
		},
	})
	log.Printf("Table name: %s", db.table)
	if err != nil {
		log.Printf("failed to get visitor count: %v", err)
		return 0, fmt.Errorf("failed to get visitor count: %v", err)
	}

	if result.Item == nil {
		return 0, nil
	}

	var visitorCount models.VisitorCount
	err = attributevalue.UnmarshalMap(result.Item, &visitorCount)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal visitor count: %v", err)
	}

	return visitorCount.Count, nil
}
