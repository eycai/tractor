package models

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

var Suits = []Suit{Spade, Diamond, Heart, Club, Joker}

type Team string

const (
	Bosses   Team = "BOSSES"
	Peasants Team = "PEASANTS"
)
