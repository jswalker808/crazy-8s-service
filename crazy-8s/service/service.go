package service

import (
	"crazy-8s/transport"
	"log"
)

type GameService struct {

}

func NewGameService() *GameService {
	return &GameService{}
}

func (service *GameService) CreateGame(request *transport.CreateGameRequest) error {
	log.Println("In the game service")

	log.Printf("Player name %v", request.PlayerName)
	return nil
}