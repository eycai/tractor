package models

type Room struct {
	ID      int
	Players []Player
	Game    Game
}

type Game struct {
	Points    []int
	Trump     int
	Vault     int // points in vault
	BankerID  int
	PlayOrder []int
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

const (
	Ace   Value = 1
	Two   Value = 2
	Three Value = 3
	Four  Value = 4
	Five  Value = 5
	Six   Value = 6
	Seven Value = 7
	Eight Value = 8
	Nine  Value = 9
	Ten   Value = 10
	Jack  Value = 11
	Queen Value = 12
	King  Value = 13
	Big   Value = 101
	Small Value = 100
)
