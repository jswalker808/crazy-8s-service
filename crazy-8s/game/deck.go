package game

import "errors"

type Deck interface {
	Draw() (*Card, error)
	Shuffle()
}

type StandardDeck struct {
	cards []Card
}

const maxCards = 52

func NewStandardDeck() *StandardDeck {
	cards := make([]Card, maxCards)



	return &StandardDeck{
		cards: cards,
	}
}

func (deck *StandardDeck) Draw() (*Card, error) {
	if len(deck.cards) == 0 {
		return nil, errors.New("deck is empty, unable to draw")
	}

	lastIdx := len(deck.cards) - 1
	top := &deck.cards[lastIdx:][0]
	deck.cards = deck.cards[:len(deck.cards) - 1]

	return top, nil
}

func (deck *StandardDeck) Shuffle() {
}

