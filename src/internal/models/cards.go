package models

import (
	"math/rand"
	"time"
)

type Card struct {
	Value   int  `json:"value"`
	Suit    Suit `json:"suit"`
	IsTrump bool `json:"isTrump"`
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

func (d *Deck) shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

type Pattern string

const (
	NOfAKind Pattern = "N_OF_A_KIND"
	Tractor  Pattern = "TRACTOR"
)

type Trick struct {
	Pattern               Pattern
	LargestCard           int
	NumCards              int
	TractorNumConsecutive int
	Suit                  Suit
	IsTrump               bool
}

func parseTrick(cards []Card) Trick {
	trick := Trick{
		NumCards: len(cards),
		Suit:     cards[0].Suit,
		IsTrump:  cards[0].IsTrump,
	}

	// for i, c := cards {

	// }
	return trick
}

func getTricks(cards [][]Card, trumpSuit Suit, trumpNumber int) ([]Trick, error) {
	tricks := make([]Trick, len(cards))
	for i, t := range cards {
		tricks[i] = parseTrick(t)
	}
	return tricks, nil
}

// sorting
type ByValue []Card

func (v ByValue) Len() int           { return len(v) }
func (v ByValue) Less(i, j int) bool { return v[i].Value < v[j].Value }
func (v ByValue) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }

// func subCards(play []Card) [][]Card {

// }

// func beats(a []Card, b []Card) bool {

// }

// func matches(a []Card, b []Card) bool {

// }
