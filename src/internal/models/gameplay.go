package models

type Room struct {
	ID     int      `json:"id"`
	Users  []string `json:"users"`
	HostID string   `json:"hostId"`
	Game   Game     `json:"Game"`
}

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
	ID       int
	Name     string
	TeamID   int
	IsBanker bool
	Level    int
	Cards    []Card
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
