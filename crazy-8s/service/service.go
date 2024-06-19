package service

import (
	gamePkg "crazy-8s/game"
	"crazy-8s/repository"
	"crazy-8s/transport"
	"log"
)

type GameService struct {
	gameRepository *repository.GameRepository
}

func NewGameService(gameRepository *repository.GameRepository) *GameService {
	return &GameService{
		gameRepository: gameRepository,
	}
}

func (service *GameService) CreateGame(connectionId string, request *transport.CreateGameRequest) (*gamePkg.Game, error) {
	log.Println("Creating new game")
	log.Printf("connectionId: %v", connectionId)

	createdGame, createGameErr := service.gameRepository.CreateGame(gamePkg.NewGame(connectionId, request.PlayerName))
	if createGameErr != nil {
		return nil, createGameErr
	}

	log.Printf("Game was successfully added to database")

	return createdGame, nil
}

func (service *GameService) JoinGame(connectionId string, request *transport.JoinGameRequest) (*gamePkg.Game, error) {
	log.Println("Joining game")
	log.Printf("connectionId: %v", connectionId)

	game, getGameErr := service.gameRepository.GetGame(request.GameId)
	if getGameErr != nil {
		return nil, getGameErr
	}

	newPlayer := gamePkg.NewPlayer(connectionId, request.PlayerName)
	addPlayerErr := game.AddPlayer(newPlayer)
	if addPlayerErr != nil {
		return nil, addPlayerErr
	}

	storePlayerErr := service.gameRepository.AddPlayer(request.GameId, newPlayer)
	if storePlayerErr != nil {
		return nil, storePlayerErr
	}

	log.Printf("game with added player: %v", game)

	return game, nil
}


func (service *GameService) LeaveGame(connectionId string) (*gamePkg.Game, error) {
	log.Println("Leaving game")
	log.Printf("connectionId: %v", connectionId)

	updatedGame,removePlayerErr := service.gameRepository.RemovePlayer(connectionId)
	if removePlayerErr != nil {
		return nil, removePlayerErr
	}

	return updatedGame, nil
}
