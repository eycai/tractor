package models

type Game struct {
	Players     []Player          `json:"players"`
	Kitty       []Card            `json:"kitty"`
	TurnUserID  string            `json:"turnUserId"`
	TrumpSuit   Suit              `json:"TrumpSuit"`
	TrumpNumber int               `json:"TrumpNumber"`
	BankerID    string            `json:"bankerID"`
	CardsInPlay map[string][]Card `json:"cardsInPlay"`
}

type Player struct {
	ID     int
	Team   Team
	Level  int
	Hand   []Card
	Points int
}

type Card struct {
	Value int
	Suit  Suit
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
