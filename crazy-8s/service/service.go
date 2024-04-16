package service

import (
	"crazy-8s/game"
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

func (service *GameService) CreateGame(connectionId string, request *transport.CreateGameRequest) error {
	log.Println("Creating new game")
	log.Printf("connectionId: %v", connectionId)

	_, err := service.gameRepository.CreateGame(game.NewGame(connectionId))
	if err != nil {
		return err
	}

	log.Printf("Game was successfully created")

	return nil
}