package transport

import "crazy-8s/game"

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

func NewGameResponse(game *game.Game, connectionId string) GameResponse {
	opponents := make([]PlayerResponse, 0)
	for _, player := range game.GetOpponents(connectionId) {
		opponents = append(opponents, NewPlayerResponse(player))
	}

	deck := make([]CardResponse, 0)
	for _, card := range game.GetDeck().GetCards() {
		deck = append(deck, NewCardResponse(card))
	}

	discardPile := make([]CardResponse, 0)
	for _, card := range game.GetDiscardPile() {
		discardPile = append(discardPile, NewCardResponse(card))
	}

	return GameResponse{
		GameId:      game.GetId(),
		MaxPoints:   game.GetMaxPoints(),
		Owner:       game.GetOwnerId(),
		Player:      NewPlayerResponse(game.GetPlayer(connectionId)),
		Opponents:   opponents,
		Deck:        deck,
		DiscardPile: discardPile,
		Status:      string(game.GetStatus()),
		CurrentTurn: game.GetCurrentTurn(),
	}
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

func NewCardResponse(card *game.Card) CardResponse {
	return CardResponse{
		Number: card.GetNumber(),
		Color:  string(card.GetColor()),
	}
}
