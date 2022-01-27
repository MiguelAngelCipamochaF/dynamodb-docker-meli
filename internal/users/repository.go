package users

import (
	"context"

	"github.com/MiguelAngelCipamochaF/dynamodb-docker-meli/internal/users/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type RepositoryDynamo interface {
	Store(ctx context.Context, model *models.User) error
	GetOne(ctx context.Context, id string) (*models.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, firstname string, lastname string, email string) error
}

type dynamoRepository struct {
	dynamo *dynamodb.DynamoDB
	table  string
}

func NewDynamoRepository(dynamo *dynamodb.DynamoDB, table string) RepositoryDynamo {
	return &dynamoRepository{
		dynamo: dynamo,
		table:  table,
	}
}

func (r *dynamoRepository) Store(ctx context.Context, model *models.User) error {
	av, err := dynamodbattribute.MarshalMap(model)

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.table),
	}

	_, err = r.dynamo.PutItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (r *dynamoRepository) GetOne(ctx context.Context, id string) (*models.User, error) {
	result, err := r.dynamo.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	return models.ItemToUser(result.Item)
}

func (r *dynamoRepository) Update(ctx context.Context, id string, firstname string, lastname string, email string) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":fn": {
				S: aws.String(firstname),
			},
			":ln": {
				S: aws.String(lastname),
			},
			":e": {
				S: aws.String(email),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName:        aws.String(r.table),
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set firstname = :fn, lastname = :ln, email = :e"),
	}
	_, err := r.dynamo.UpdateItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (r *dynamoRepository) Delete(ctx context.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(r.table),
	}
	_, err := r.dynamo.DeleteItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
