package game

const startingHandLength = 10

type Player struct {
	id string
	name string
	hand []Card
	points int
}

func NewPlayer(id string, name string) *Player {
	return &Player{
		id: id,
		name: name,
		hand: make([]Card, startingHandLength),
		points: 0,
	}
}