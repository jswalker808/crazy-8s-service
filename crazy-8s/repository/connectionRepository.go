package repository

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ConnectionRepository struct {
	dynamoDbClient *dynamodb.Client
}

func NewConnectionRepository (dynamoDbClient *dynamodb.Client) *ConnectionRepository {
	return &ConnectionRepository{
		dynamoDbClient: dynamoDbClient,
	}
}

func (repository *ConnectionRepository) AddConnection(connectionId string, gameId string) error {
	log.Printf("Adding connection %v to game %v", connectionId, gameId)

	attributeMap := make(map[string]types.AttributeValue)
	attributeMap["connectionId"] = &types.AttributeValueMemberS{Value: connectionId}
	attributeMap["gameId"] = &types.AttributeValueMemberS{Value: gameId}

	_, putItemErr := repository.dynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("CONNECTIONS_TABLE_NAME")),
		Item: attributeMap,
	})
	if putItemErr != nil {
		return putItemErr
	}

	return nil
}
