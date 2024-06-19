package api

import (
	"context"
	"crazy-8s/global"
	"crazy-8s/notification"
	"crazy-8s/service"
	"crazy-8s/transport"
	"fmt"
)

type Router struct {
	gameService *service.GameService
	notifier    notification.Notifier
}

func NewRouter(gameService *service.GameService, notifier notification.Notifier) *Router {
	return &Router{
		gameService: gameService,
		notifier:    notifier,
	}
}

func (router *Router) HandleDisconnect(connectionId string) error {
	updatedGame, err := router.gameService.LeaveGame(connectionId)
	if err != nil {
		return err
	}

	router.notifier.SendAll(transport.NewGameResponseMap(updatedGame))

	return nil
}

func (router *Router) HandleRequest(ctx context.Context, action string, request transport.Request) error {
	switch action {
		case "create_game":
			return router.handleCreateGame(ctx, request)
		case "join_game":
			return router.handleJoinGame(ctx, request)
		default:
			return fmt.Errorf("unsupported game action: %v", action)
	}
}

func (router *Router) handleCreateGame(ctx context.Context, request transport.Request) error {
	gameRequest, ok := request.(*transport.CreateGameRequest)
	if !ok {
		return fmt.Errorf("CreateGameRequest is required to create a new game")
	}

	connectionId := ctx.Value(global.ConnectionIdCtxKey{}).(string)

	createdGame, err := router.gameService.CreateGame(connectionId, gameRequest)

	if notifyErr := router.notifier.Send(connectionId, transport.NewGameResponse(createdGame, connectionId)); notifyErr.Error != nil {
		return notifyErr.Error
	}

	return err
}

func (router *Router) handleJoinGame(ctx context.Context, request transport.Request) error {
	gameRequest, ok := request.(*transport.JoinGameRequest)
	if !ok {
		return fmt.Errorf("JoinGameRequest is required to join an existing game")
	}
	
	game, err := router.gameService.JoinGame(ctx.Value(global.ConnectionIdCtxKey{}).(string), gameRequest)
	if err != nil {
		return err
	}

	router.notifier.SendAll(transport.NewGameResponseMap(game))

	return nil
}

func (router *Router) Notifier() notification.Notifier {
	return router.notifier
}

