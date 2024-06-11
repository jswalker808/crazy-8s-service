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

	gameAttributeValue, marshalErr := attributevalue.MarshalMap(gameStore)
	if marshalErr != nil {
		return nil, marshalErr
	}

	connectionAttributeValue := map[string]types.AttributeValue{
		"gameId": &types.AttributeValueMemberS{Value: game.GetId()},
		"connectionId": &types.AttributeValueMemberS{Value: game.GetOwnerId()},
	}

	_, transactionErr := repository.dynamoDbClient.TransactWriteItems(context.TODO(), &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName: aws.String(getGameConnectionsTableName()),
					Item:      connectionAttributeValue,
				},
			},
			{
				Put: &types.Put{
					TableName: aws.String(getGameTableName()),
					Item:      gameAttributeValue,
				},
			},
		},
	})

	if transactionErr != nil {
		return nil, transactionErr
	}

	return game, nil
}

func (repository *GameRepository) GetGame(gameId string) (*gamePkg.Game, error) {
	getItemOutput, getItemErr := repository.dynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"gameId": &types.AttributeValueMemberS{Value: gameId},
		},
		TableName: aws.String(getGameTableName()),
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

	connectionAttributeValue := map[string]types.AttributeValue{
		"gameId": &types.AttributeValueMemberS{Value: gameId},
		"connectionId": &types.AttributeValueMemberS{Value: player.GetId()},
	}

	_, transactionErr := repository.dynamoDbClient.TransactWriteItems(context.TODO(), &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName: aws.String(getGameConnectionsTableName()),
					Item:      connectionAttributeValue,
				},
			},
			{
				Update: &types.Update{
					Key: map[string]types.AttributeValue{
						"gameId": &types.AttributeValueMemberS{Value: gameId},
					},
					TableName: aws.String(getGameTableName()),
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
				},
			},
		},
	})

	if transactionErr != nil {
		return transactionErr
	}

	return nil
}

func (repository *GameRepository) RemovePlayer(connectionId string) (*gamePkg.Game, error) {
	connectionAttributeValue := map[string]types.AttributeValue{
		"connectionId": &types.AttributeValueMemberS{Value: connectionId},
	}

	deleteConnectionOutput, delectConnectionErr := repository.dynamoDbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key: connectionAttributeValue,
		TableName: aws.String(getGameConnectionsTableName()),
		ReturnValues: types.ReturnValueAllOld,
	})

	if (delectConnectionErr != nil) {
		return nil, delectConnectionErr
	}

	// Do nothing if connection is not tied to any game
	if (deleteConnectionOutput.Attributes == nil) {
		return nil, nil
	}

	gameId := deleteConnectionOutput.Attributes["gameId"].(*types.AttributeValueMemberS).Value
	game, getGameErr := repository.GetGame(gameId)
	if getGameErr != nil {
		return nil, getGameErr
	}

	game.RemovePlayer(connectionId)

	gameStore := NewGameStore(game)

	playersAttributeValue, marshalErr := attributevalue.MarshalMap(map[string]interface{}{
		":updatedPlayers": gameStore.Players,
		":updatedDeck": gameStore.Deck,
	})
	if marshalErr != nil {
		return nil, marshalErr
	}

	_, updateGameErr := repository.dynamoDbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"gameId": &types.AttributeValueMemberS{Value: gameId},
		},
		TableName: aws.String(getGameTableName()),
		UpdateExpression: aws.String("Set Players = :updatedPlayers, Deck = :updatedDeck"),
		ExpressionAttributeValues: playersAttributeValue,
	})

	if updateGameErr != nil {
		return nil, updateGameErr
	}

	return game, nil
}

func getGameTableName() string {
	return os.Getenv("TABLE_NAME")
}

func getGameConnectionsTableName() string {
	return os.Getenv("CONNECTIONS_TABLE_NAME")
}
