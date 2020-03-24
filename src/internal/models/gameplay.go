package models

type Game struct {
	Players     []*Player         `json:"players"`
	Turn        string            `json:"turn"`
	TrumpSuit   Suit              `json:"trumpSuit"`
	TrumpNumber int               `json:"trumpNumber"`
	Banker      string            `json:"banker"`
	CardsInPlay map[string][]Card `json:"cardsInPlay"`
}

type Player struct {
	Username string `json:"username"`
	Team     Team   `json:"team"`
	Level    int    `json:"level"`
	Points   int    `json:"points"`
}

type Card struct {
	Value int  `json:"value"`
	Suit  Suit `json:"suit"`
}

type Deck struct {
	Cards []Card
}

// suit enum
type Suit string

const (
	Spade   Suit = "SPADE"
	Diamond Suit = "DIAMOND"
	Heart   Suit = "HEART"
	Club    Suit = "CLUB"
	Joker   Suit = "JOKER"
)

// card values
type Value int

type Team string

const (
	Bosses   Team = "BOSSES"
	Peasants Team = "PEASANTS"
)
