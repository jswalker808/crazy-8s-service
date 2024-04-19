package transport

import "crazy-8s/game"

type CreateGameResponse struct {
	GameId string
	MaxPoints int
	Owner CreateGameOwnerResponse
}

type CreateGameOwnerResponse struct {
	Id string
	Name string
}

func NewCreateGameResponse(game *game.Game) CreateGameResponse {
	return CreateGameResponse{
		GameId: game.GetId(),
		MaxPoints: game.GetMaxPoints(),
		Owner: NewCreateGameOwnerResponse(game.GetOwner()),
	}
}

func NewCreateGameOwnerResponse(player *game.Player) CreateGameOwnerResponse {
	return CreateGameOwnerResponse{
		Id: player.GetId(),
		Name: player.GetName(),
	}
}

