package game

import "fmt"

type Color string

const (
	RED Color = "red"
	BLUE = "blue"
	GREEN = "green"
	YELLOW = "yellow"
)

var Colors = []Color {
	RED,
	BLUE,
	GREEN,
	YELLOW,
}

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

func (card *Card) String() string {
	return fmt.Sprintf("Card[color=%v, number=%v]", card.color, card.number)
}

func (card *Card) GetColor() Color {
	return card.color
}

func (card *Card) GetNumber() int {
	return card.number
}



