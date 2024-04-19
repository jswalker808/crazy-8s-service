package api

import (
	"context"
	"crazy-8s/global"
	"crazy-8s/service"
	"crazy-8s/transport"
	"fmt"
)

type Router struct {
	gameService *service.GameService
}

func NewRouter(gameService *service.GameService) *Router {
	return &Router{
		gameService: gameService,
	}
}

func (router *Router) HandleRequest(ctx context.Context, action string, request transport.Request) error {
	switch action {
		case "create_game":
			return router.handleCreateGame(ctx, request)
		default:
			return fmt.Errorf("unsupported game action: %v", action)
	}
}

func (router *Router) handleCreateGame(ctx context.Context, request transport.Request) error {
	gameRequest, ok := request.(*transport.CreateGameRequest)
	if !ok {
		return fmt.Errorf("CreateGameRequest is required to create a new game")
	}
	return router.gameService.CreateGame(ctx.Value(global.ConnectionIdCtxKey{}).(string), gameRequest)
}

func (router *Router) GameService() *service.GameService {
	return router.gameService
}

