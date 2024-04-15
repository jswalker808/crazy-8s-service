package game

type Game struct {
	name string
}

func NewGame() *Game {
	return &Game{ name: "string" }
}