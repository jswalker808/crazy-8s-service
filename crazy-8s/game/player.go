package game

import (
	"fmt"
	"strings"
)

const startingHandLength = 10

type Player struct {
	id string
	name string
	hand []*Card
	points int
}

func NewPlayer(id string, name string) *Player {
	return &Player{
		id: id,
		name: name,
		hand: make([]*Card, 0, startingHandLength),
		points: 0,
	}
}

func NewPlayerFromExisting(id string, name string, hand []*Card, points int) *Player {
	return &Player{
		id: id,
		name: name,
		hand: hand,
		points: points,
	}
}

func (player *Player) String() string {
	var sb strings.Builder
	for _, card := range player.hand {
		sb.WriteString(fmt.Sprint(card))
		sb.WriteString(",")
	}
	return fmt.Sprintf("Player[id=%v, name=%v, hand=%v, points=%v]", player.id, player.name, sb.String(), player.points)
}

func (player *Player) GetId() string {
	return player.id
}

func (player *Player) GetName() string {
	return player.name
}

func (player *Player) GetHand() []*Card {
	return player.hand
}

func (player *Player) GetPoints() int {
	return player.points
}