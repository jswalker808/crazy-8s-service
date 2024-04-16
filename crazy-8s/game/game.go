package game

const maxPlayers = 4

type Game struct {
	id string
	ownerId string
	maxPoints int
	players []Player
	deck Deck
	discard []Card
}

func NewGame(ownerId string) *Game {
	return &Game{
		id: "slkajsdf",
		ownerId: ownerId,
		maxPoints: 100,
		players: make([]Player, maxPlayers),
		deck: NewStandardDeck(),
		discard: make([]Card, 0),
	}
}