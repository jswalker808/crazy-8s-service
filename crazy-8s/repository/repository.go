package repository

import (
	"context"
	gamePkg "crazy-8s/game"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type GameRepository struct  {
	dynamoDbClient *dynamodb.Client
}

func NewGameRepository(dynamoDbClient *dynamodb.Client) *GameRepository {
	return &GameRepository{
		dynamoDbClient: dynamoDbClient,
	}
}

func (repository *GameRepository) CreateGame(game *gamePkg.Game) (*gamePkg.Game, error) {
	log.Println(game)

	gameStore := NewGameStore(game)

	attributeValue, marshalErr := attributevalue.MarshalMap(gameStore)
	if marshalErr != nil {
		return nil, marshalErr
	}

	log.Println(os.Getenv("TABLE_NAME"))
	log.Println(attributeValue)

	_, putItemErr := repository.dynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item: attributeValue,
	})
	if putItemErr != nil {
		return nil, putItemErr
	}

	return game, nil
}