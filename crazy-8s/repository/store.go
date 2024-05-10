package repository

import gamePkg "crazy-8s/game"

type GameStore struct {
	Id string `dynamodbav:"gameId"`
	OwnerId string
	MaxPoints int
	Players []PlayerStore
	Deck []CardStore
	DiscardPile []CardStore
	Status string
	CurrentTurn string
} 

type PlayerStore struct {
	Id string
	Name string
	Hand []CardStore
	Points int
}

type CardStore struct {
	Color gamePkg.Color
	Number int
}

func NewGameStore(game *gamePkg.Game) GameStore {
	players := make([]PlayerStore, 0)
	for _, player := range game.GetPlayers() {
		players = append(players, NewPlayerStore(player))
	}

	deckCards := make([]CardStore, 0)
	for _, card := range game.GetDeck().GetCards() {
		deckCards = append(deckCards, NewCardStore(card))
	}

	discardPile := make([]CardStore, 0)
	for _, card := range game.GetDiscardPile() {
		discardPile = append(discardPile, NewCardStore(card))
	}

	return GameStore{
		Id: game.GetId(),
		OwnerId: game.GetOwnerId(),
		MaxPoints: game.GetMaxPoints(),
		Players: players,
		Deck: deckCards,
		DiscardPile: discardPile,
		Status: string(game.GetStatus()),
		CurrentTurn: game.GetCurrentTurn(),
	}
}

func NewGameFromStore(gameStore GameStore) (*gamePkg.Game , error) {
	players := make([]*gamePkg.Player, 0)
	for _, playerStore := range gameStore.Players {
		hand := make([]*gamePkg.Card, 0)
		for _, cardStore := range playerStore.Hand {
			hand = append(hand, gamePkg.NewCard(cardStore.Color, cardStore.Number))
		}
		players = append(players, gamePkg.NewPlayerFromExisting(playerStore.Id, playerStore.Name, hand, playerStore.Points))
	}

	deck := make([]*gamePkg.Card, 0)
	for _, cardStore := range gameStore.Deck {
		deck = append(deck, gamePkg.NewCard(cardStore.Color, cardStore.Number))
	}

	discardPile := make([]*gamePkg.Card, 0)
	for _, cardStore := range gameStore.DiscardPile {
		discardPile = append(discardPile, gamePkg.NewCard(cardStore.Color, cardStore.Number))
	}

	status, parseErr := gamePkg.ParseStatus(gameStore.Status)
	if parseErr != nil {
		return nil, parseErr
	}

	return gamePkg.NewGameFromExisting(
		gameStore.Id,
		gameStore.OwnerId,
		gameStore.MaxPoints,
		players,
		deck,
		discardPile,
		status,
		gameStore.CurrentTurn,
	), nil
}

func NewPlayerStore(player *gamePkg.Player) PlayerStore {

	hand := make([]CardStore, 0)
	for _, card := range player.GetHand() {
		hand = append(hand, NewCardStore(card))
	}

	return PlayerStore{
		Id: player.GetId(),
		Name: player.GetName(),
		Hand: hand,
		Points: player.GetPoints(),
	}
}

func NewCardStore(card *gamePkg.Card) CardStore {
	return CardStore{
		Color: card.GetColor(),
		Number: card.GetNumber(),
	}
}