package repository

import (
	"context"
	gamePkg "crazy-8s/game"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var tableName = os.Getenv("TABLE_NAME")

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

	log.Println(attributeValue)

	_, putItemErr := repository.dynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: attributeValue,
	})
	if putItemErr != nil {
		return nil, putItemErr
	}

	return game, nil
}

func (repository *GameRepository) GetGame(gameId string) (*gamePkg.Game, error) {
	getItemOutput, getItemErr := repository.dynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"gameId": &types.AttributeValueMemberS{Value: gameId},
		},
		TableName: aws.String(tableName),
	});

	if getItemErr != nil {
		return nil, getItemErr
	}

	var gameStore GameStore
	unmarshalErr := attributevalue.UnmarshalMap(getItemOutput.Item, &gameStore)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	game, mappingErr := NewGameFromStore(gameStore)
	if mappingErr != nil {
		return nil, mappingErr
	}

	return game, nil
}

func (repository *GameRepository) AddPlayer(gameId string, player *gamePkg.Player) error {
	log.Println(player)

	playerStore := NewPlayerStore(player)

	attributeValue, marshalErr := attributevalue.MarshalMap(playerStore)
	if marshalErr != nil {
		return marshalErr
	}

	log.Println(attributeValue)

	_, updateItemErr := repository.dynamoDbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"gameId": &types.AttributeValueMemberS{Value: gameId},
		},
		TableName: aws.String(tableName),
		UpdateExpression: aws.String("SET Players = list_append(Players, :newPlayers)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":newPlayers": &types.AttributeValueMemberL{
				Value: []types.AttributeValue{
                    &types.AttributeValueMemberM{
                        Value: attributeValue,
                    },
                },
			},
		},
		ReturnValues: types.ReturnValueNone,
	})
	if updateItemErr != nil {
		return updateItemErr
	}

	return nil
}
