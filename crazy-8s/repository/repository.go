package repository

import (
	gamePkg "crazy-8s/game"

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
	return game, nil
}