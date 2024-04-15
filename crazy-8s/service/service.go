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

func (service *GameService) CreateGame(request *transport.CreateGameRequest) error {
	log.Println("In the game service")

	createdGame, err := service.gameRepository.CreateGame(game.NewGame())
	if err != nil {
		return err
	}

	log.Printf("Game was successfully created: %v", createdGame)

	return nil
}