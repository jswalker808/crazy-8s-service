package service

import (
	"crazy-8s/game"
	"crazy-8s/notification"
	"crazy-8s/repository"
	"crazy-8s/transport"
	"encoding/json"
	"log"
)

type GameService struct {
	gameRepository *repository.GameRepository
	notifier       *notification.ApiGatewayNotifier
}

func NewGameService(gameRepository *repository.GameRepository, notifier *notification.ApiGatewayNotifier) *GameService {
	return &GameService{
		gameRepository: gameRepository,
		notifier:       notifier,
	}
}

func (service *GameService) CreateGame(connectionId string, request *transport.CreateGameRequest) error {
	log.Println("Creating new game")
	log.Printf("connectionId: %v", connectionId)

	createdGame, createGameErr := service.gameRepository.CreateGame(game.NewGame(connectionId, request.PlayerName))
	if createGameErr != nil {
		return createGameErr
	}

	log.Printf("Game was successfully added to database")

	createdGameBytes, jsonErr := json.Marshal(transport.NewGameResponse(createdGame, connectionId))
	if jsonErr != nil {
		panic("Unable to marshal created game to JSON")
	}

	log.Println(createdGameBytes)

	if notifyErr := service.notifier.Send(connectionId, createdGameBytes); notifyErr != nil {
		return notifyErr
	}

	return nil
}

func (service *GameService) Notifier() *notification.ApiGatewayNotifier {
	return service.notifier
}
