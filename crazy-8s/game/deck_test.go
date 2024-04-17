package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStandardDeck(t *testing.T) {
	assert := assert.New(t)

	standardDeck := NewStandardDeck()

	assert.NotNil(standardDeck)
	assert.Equal(52, len(standardDeck.cards))
}

func TestStandardDeckDraw(t *testing.T) {
	assert := assert.New(t)

	standardDeck := NewStandardDeck()

	assert.Equal(52, len(standardDeck.cards))

	top, err := standardDeck.Draw()
	assert.Nil(err)
	assert.NotNil(top)
	assert.Equal(51, len(standardDeck.cards))
}