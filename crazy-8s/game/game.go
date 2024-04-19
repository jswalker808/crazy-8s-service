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
	discardPile []*Card
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
		discardPile: make([]*Card, 0),
	}
}

func (game *Game) String() string {
	return fmt.Sprintf("Game[id=%v, ownerId=%v, maxPoints=%v, players=%v]", game.id, game.ownerId, game.maxPoints, len(game.players))
}

func (game *Game) GetId() string {
	return game.id
}

func (game *Game) GetOwnerId() string {
	return game.ownerId
}

func (game *Game) GetMaxPoints() int {
	return game.maxPoints
}

func (game *Game) GetPlayers() []*Player {
	return game.players
} 

func (game *Game) GetDeck() Deck {
	return game.deck
}

func (game *Game) GetDiscardPile() []*Card {
	return game.discardPile
}

func (game *Game) GetOwner() *Player {
	for _, player := range game.players {
		if player.id == game.ownerId {
			return player
		}
	}
	panic("Invalid state: Unable to find owner for game")
}