package service

import (
	gamePkg "crazy-8s/game"
	"crazy-8s/notification"
	"crazy-8s/repository"
	"crazy-8s/transport"
	"encoding/json"
	"log"
)

type GameService struct {
	gameRepository *repository.GameRepository
	notifier       notification.Notifier
}

func NewGameService(
	gameRepository *repository.GameRepository,
	notifier notification.Notifier) *GameService {
	return &GameService{
		gameRepository: gameRepository,
		notifier:       notifier,
	}
}

func (service *GameService) Notifier() notification.Notifier {
	return service.notifier
}

func (service *GameService) CreateGame(connectionId string, request *transport.CreateGameRequest) (*gamePkg.Game, error) {
	log.Println("Creating new game")
	log.Printf("connectionId: %v", connectionId)

	createdGame, createGameErr := service.gameRepository.CreateGame(gamePkg.NewGame(connectionId, request.PlayerName))
	if createGameErr != nil {
		return nil, createGameErr
	}

	log.Printf("Game was successfully added to database")

	createdGameBytes, jsonErr := json.Marshal(transport.NewGameResponse(createdGame, connectionId))
	if jsonErr != nil {
		return nil, jsonErr
	}

	if notifyErr := service.notifier.Send(connectionId, createdGameBytes); notifyErr != nil {
		return nil, notifyErr
	}

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

	gameResponseMap := transport.NewGameResponseMap(game)

	gameBytesMap := make(map[string][]byte)
	for connectionId, gameResponse := range gameResponseMap {
		gameBytes, jsonErr := json.Marshal(gameResponse)
		if jsonErr != nil {
			return nil, jsonErr
		}
		gameBytesMap[connectionId] = gameBytes
	}

	notificationErrors := service.notifier.SendAll(gameBytesMap)
	for _, err := range notificationErrors {
		log.Printf("Connection %v ran into an error: %v", err.ConnectionId, err.Error)
	}

	return game, nil
}

