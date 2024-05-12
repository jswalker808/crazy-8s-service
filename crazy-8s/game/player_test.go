package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlayer(t *testing.T) {
	assert := assert.New(t)

	givenPlayerId := "playerId"
	givenPlayerName := "playerName"

	actual := NewPlayer(givenPlayerId, givenPlayerName)

	assert.Equal(givenPlayerId, actual.GetId())
	assert.Equal(givenPlayerName, actual.GetName())
	assert.Empty(actual.GetHand())
}
