package game

type Color string

const (
	RED = iota
	BLUE
	GREEN
	YELLOW
)

type Card struct {
	color Color
	number int
}

func NewCard(color Color, number int) *Card {
	return &Card{
		color: color,
		number: number,
	}
}


