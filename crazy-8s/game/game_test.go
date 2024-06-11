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

func TestGetOwner(t *testing.T) {
	assert := assert.New(t)

	givenOwnerId := "ownerId"
	givenOwnerName := "TestUser"

	actual := NewGame(givenOwnerId, givenOwnerName)
	
	owner := actual.GetOwner()
	assert.NotNil(owner)
	assert.Equal(givenOwnerId, owner.GetId())
	assert.Equal(givenOwnerName, owner.GetName())
}

func TestGetPlayer(t *testing.T) {
	assert := assert.New(t)

	givenOwnerId := "ownerId"
	givenOwnerName := "TestUser"
	givenGame := NewGame(givenOwnerId, givenOwnerName)

	var tests = []struct {
		name string
		givenPlayerId string
		playerFound bool
	} {
		{"Existing Player", givenOwnerId, true},
		{"Non-existent player", "bogusPlayerId", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := givenGame.GetPlayer(test.givenPlayerId)
			assert.Equal(test.playerFound, actual != nil)
		})
	}
}

func TestAddPlayer(t *testing.T) {
	assert := assert.New(t)

	givenOwnerId := "ownerId"
	givenOwnerName := "TestUser"

	actual := NewGame(givenOwnerId, givenOwnerName)
	assert.Equal(1, len(actual.GetPlayers()))

	err := actual.AddPlayer(NewPlayer("newPlayerId", "newPlayerName"))
	assert.Equal(2, len(actual.GetPlayers()))
	assert.NoError(err)
}

func TestAddPlayerMax(t *testing.T) {
	assert := assert.New(t)

	givenOwnerId := "ownerId"
	givenOwnerName := "TestUser"

	actual := NewGame(givenOwnerId, givenOwnerName)
	assert.Equal(1, len(actual.GetPlayers()))

	for i := 1; i < MaxPlayers; i++ {
		assert.NoError(actual.AddPlayer(NewPlayer("playerId", "playerName")))
		assert.Equal(i + 1, len(actual.GetPlayers()))
	}

	assert.Error(actual.AddPlayer(NewPlayer("playerId", "playerName")))
}

func TestRemovePlayer(t *testing.T) {
	assert := assert.New(t)

	givenOwnerId := "ownerId"
	givenOwnerName := "TestUser"
	givenGame := NewGame(givenOwnerId, givenOwnerName)
	assert.Equal(1, len(givenGame.GetPlayers()))
	assert.Equal(52, len(givenGame.GetDeck().GetCards()))

	newCards := make([]*Card, 0)
	newCards = append(newCards, NewCard(RED, 0))
	newCards = append(newCards, NewCard(BLUE, 1))
	newPlayer := NewPlayer("newPlayerId", "newPlayerName")
	newPlayer.AddCards(newCards)
	assert.Equal(2, len(newPlayer.GetHand()))

	givenGame.AddPlayer(newPlayer)
	assert.Equal(2, len(givenGame.GetPlayers()))

	givenGame.RemovePlayer("newPlayerId")
	assert.Equal(1, len(givenGame.GetPlayers()))
	assert.Equal(54, len(givenGame.GetDeck().GetCards()))
}

