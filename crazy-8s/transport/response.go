package transport

import (
	"crazy-8s/game"
	"encoding/json"
	"fmt"
	"log"
)

type CardResponse struct {
	Number int    `json:"number"`
	Color  string `json:"color"`
}

type PlayerResponse struct {
	Id     string         `json:"id"`
	Name   string         `json:"name"`
	Points int            `json:"points"`
	Hand   []CardResponse `json:"hand"`
}

type GameResponse struct {
	GameId      string           `json:"gameId"`
	MaxPoints   int              `json:"maxPoints"`
	Owner       string           `json:"owner"`
	Player      PlayerResponse   `json:"player"`
	Opponents   []PlayerResponse `json:"opponents"`
	Deck        []CardResponse   `json:"deck"`
	DiscardPile []CardResponse   `json:"discardPile"`
	CurrentTurn string           `json:"currentTurn"`
	Status      string           `json:"status"`
}

type NewGameResponseOption = func(*GameResponse)

func NewGameResponseMap(game *game.Game) map[string][]byte {
	responseMap := make(map[string][]byte)

	playerResponseMap := make(map[string]PlayerResponse, 0)
	for _, player := range game.GetPlayers() {
		playerResponseMap[player.GetId()] = NewPlayerResponse(player)
	}

	deck := NewCardResponses(game.GetDeck().GetCards())
	discardPile := NewCardResponses(game.GetDiscardPile())

	for _, player := range game.GetPlayers() {
		playerResponse := playerResponseMap[player.GetId()]
		opponents := NewPlayerResponses(game.GetOpponents(player.GetId()))
		responseMap[player.GetId()] = NewGameResponse(
			game,
			player.GetId(),
			WithPlayer(playerResponse),
			WithOpponents(opponents),
			WithDeck(deck),
			WithDiscardPile(discardPile),
		)
	}

	fmt.Printf("responseMap: %v\n", responseMap)

	return responseMap
}

func NewGameResponse(game *game.Game, connectionId string, options ...NewGameResponseOption) []byte {
	gameResponse := &GameResponse{
		GameId:      game.GetId(),
		MaxPoints:   game.GetMaxPoints(),
		Owner:       game.GetOwnerId(),
		Status:      string(game.GetStatus()),
		CurrentTurn: game.GetCurrentTurn(),
	}

	for _, option := range options {
		option(gameResponse)
	}

	if gameResponse.Player.Id == "" {
		gameResponse.Player = NewPlayerResponse(game.GetPlayer(connectionId))
	}

	if gameResponse.Opponents == nil {
		gameResponse.Opponents = NewPlayerResponses(game.GetOpponents(connectionId))
	}

	if gameResponse.Deck == nil {
		gameResponse.Deck = NewCardResponses(game.GetDeck().GetCards())
	}

	if gameResponse.DiscardPile == nil {
		gameResponse.DiscardPile = NewCardResponses(game.GetDiscardPile())
	}

	gameResponseBytes, jsonErr := json.Marshal(gameResponse)
	if jsonErr != nil {
		panic(jsonErr)
	}

	unmarshalErr := json.Unmarshal(gameResponseBytes, &gameResponse)
	if unmarshalErr != nil {
		panic(unmarshalErr)
	}
	log.Printf("unamrshaled gameResponse: %v", gameResponse)

	return gameResponseBytes
}

func WithPlayer(player PlayerResponse) NewGameResponseOption {
	return func(gameResponse *GameResponse) {
		gameResponse.Player = player
	}
}

func WithOpponents(opponents []PlayerResponse) NewGameResponseOption {
	return func(gameResponse *GameResponse) {
		gameResponse.Opponents = opponents
	}
}

func WithDeck(deck []CardResponse) NewGameResponseOption {
	return func(gameResponse *GameResponse) {
		gameResponse.Deck = deck
	}
}

func WithDiscardPile(discardPile []CardResponse) NewGameResponseOption {
	return func(gameResponse *GameResponse) {
		gameResponse.DiscardPile = discardPile
	}
}

func NewPlayerResponses(players []*game.Player) []PlayerResponse {
	playerResponses := make([]PlayerResponse, 0)
	for _, player := range players {
		playerResponses = append(playerResponses, NewPlayerResponse(player))
	}
	return playerResponses
}

func NewPlayerResponse(player *game.Player) PlayerResponse {
	playerHand := make([]CardResponse, 0)
	for _, card := range player.GetHand() {
		playerHand = append(playerHand, NewCardResponse(card))
	}

	return PlayerResponse{
		Id:     player.GetId(),
		Name:   player.GetName(),
		Points: player.GetPoints(),
		Hand:   playerHand,
	}
}

func NewCardResponses(cards []*game.Card) []CardResponse {
	cardResponses := make([]CardResponse, 0)
	for _, card := range cards {
		cardResponses = append(cardResponses, NewCardResponse(card))
	}
	return cardResponses
}

func NewCardResponse(card *game.Card) CardResponse {
	return CardResponse{
		Number: card.GetNumber(),
		Color:  string(card.GetColor()),
	}
}
