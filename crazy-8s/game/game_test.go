package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStatus(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		name string
		given string
		expected  Status
	} {
		{"Pending status", "pending", PENDING},
		{"In progress status", "in_progress", IN_PROGRESS},
		{"Invalid existing status", "In Progress", ""},
		{"Non-existent status", "Bogus", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, _ := ParseStatus(test.given)
			assert.Equal(test.expected, actual, "Got %s, want %s", actual, test.expected)
		})
	}
}

func TestNewGame(t *testing.T) {
	assert := assert.New(t)

	givenOwnerId := "ownerId"
	givenOwnerName := "TestUser"
	
	actual := NewGame(givenOwnerId, givenOwnerName)

	assert.NotNil(actual.GetId())
	assert.Equal(givenOwnerId, actual.GetOwnerId())
	assert.Equal(100, actual.GetMaxPoints())
	assert.Equal(1, len(actual.GetPlayers()))
	assert.IsType(&StandardDeck{}, actual.GetDeck())
	assert.Empty(actual.discardPile)
	assert.Equal(PENDING, actual.GetStatus())
	assert.Empty(actual.GetCurrentTurn())
}