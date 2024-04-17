package game

import (
	"fmt"

	"github.com/google/uuid"
)

const maxPlayers = 4

type Game struct {
	id string
	ownerId string
	maxPoints int
	players []*Player
	deck Deck
	discard []*Card
}

func NewGame(ownerId string, ownerName string) *Game {
	owner := NewPlayer(ownerId, ownerName)
	players := make([]*Player, 0, maxPlayers)
	players = append(players, owner)

	return &Game{
		id: uuid.NewString(),
		ownerId: owner.id,
		maxPoints: 100,
		players: players,
		deck: NewStandardDeck(),
		discard: make([]*Card, 0),
	}
}

func (game *Game) String() string {
	return fmt.Sprintf("Game[id=%v, ownerId=%v, maxPoints=%v, players=%v]", game.id, game.ownerId, game.maxPoints, len(game.players))
}