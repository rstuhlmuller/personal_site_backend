package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rstuhlmuller/personal_site_backend/internal/models"
)

type DynamoDBInterface interface {
	CreateProject(ctx context.Context, project *models.Project) error
	GetProject(ctx context.Context, id string) (*models.Project, error)
	ListProjects(ctx context.Context) ([]models.Project, error)
	UpdateProject(ctx context.Context, project *models.Project) error
}

type DynamoDB struct {
	client *dynamodb.Client
	table  string
}

func NewDynamoDB() *DynamoDB {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	return &DynamoDB{
		client: dynamodb.NewFromConfig(cfg),
		table:  "Personal-USW2-DDB-Personal-Site", // Make sure this matches your actual table name
	}
}

func (db *DynamoDB) CreateProject(ctx context.Context, project *models.Project) error {
	item, err := attributevalue.MarshalMap(project)
	if err != nil {
		return fmt.Errorf("failed to marshal project: %v", err)
	}

	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(db.table),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to create project: %v", err)
	}

	return nil
}

func (db *DynamoDB) GetProject(ctx context.Context, id string) (*models.Project, error) {
	result, err := db.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(db.table),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %v", err)
	}

	if result.Item == nil {
		return nil, nil // Project not found
	}

	var project models.Project
	err = attributevalue.UnmarshalMap(result.Item, &project)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal project: %v", err)
	}

	return &project, nil
}

func (db *DynamoDB) ListProjects(ctx context.Context) ([]models.Project, error) {
	result, err := db.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(db.table),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %v", err)
	}

	var projects []models.Project
	err = attributevalue.UnmarshalListOfMaps(result.Items, &projects)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal projects: %v", err)
	}

	return projects, nil
}

func (db *DynamoDB) UpdateProject(ctx context.Context, project *models.Project) error {
	item, err := attributevalue.MarshalMap(project)
	if err != nil {
		return fmt.Errorf("failed to marshal project: %v", err)
	}

	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(db.table),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to update project: %v", err)
	}

	return nil
}
