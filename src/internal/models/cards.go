package models

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Card struct {
	Value         int  `json:"value"`
	Suit          Suit `json:"suit"`
	IsTrumpSuit   bool `json:"isTrumpSuit"`
	IsTrumpNumber bool `json:"isTrumpNumber"`
	// private
	gameValue int
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
	LargestCard           Card
	NumCards              int
	TractorNumConsecutive int
	Suit                  Suit
	IsTrump               bool
}

func parseTrick(cards []Card) (Trick, error) {
	trick := Trick{
		Pattern:     NOfAKind,
		NumCards:    len(cards),
		Suit:        cards[0].Suit,
		IsTrump:     cards[0].IsTrump(),
		LargestCard: cards[0],
	}

	numConsecutive := 1
	for i := 0; i < len(cards)-1; i++ {
		if cards[i].Value != cards[i+1].Value {
			if trick.TractorNumConsecutive != 0 && trick.TractorNumConsecutive != numConsecutive {
				return trick, fmt.Errorf("invalid play")
			} else if !IsConsecutive(cards[i], cards[i+1]) {
				return trick, fmt.Errorf("invalid play")
			}
			trick.Pattern = Tractor
			trick.TractorNumConsecutive = numConsecutive
			numConsecutive = 1
		} else {
			numConsecutive++
		}
	}
	return trick, nil
}

// assumes same suit or both trump
func IsConsecutive(a Card, b Card) bool {
	if a.IsTrumpNumber {
		return IsConsecutiveTrumpNumber(a, b)
	} else if b.IsTrumpNumber {
		return IsConsecutiveTrumpNumber(b, a)
	} else if a.Suit == b.Suit {
		return math.Abs(float64(a.Value-b.Value)) == 1
	} else {
		return false
	}
}

func IsConsecutiveTrumpNumber(a Card, b Card) bool {
	if a.IsTrumpSuit {
		return (b.IsTrumpNumber && !b.IsTrumpSuit) || (b.Suit == Joker && b.Value == 1)
	} else {
		return (b.IsTrumpNumber && b.IsTrumpSuit) || (b.IsTrumpSuit && b.Value == 1)
	}
}

func getTricks(cards [][]Card, trumpSuit Suit, trumpNumber int) ([]Trick, error) {
	tricks := make([]Trick, len(cards))
	for i, t := range cards {
		trick, err := parseTrick(t)
		if err != nil {
			return tricks, err
		}
		tricks[i] = trick
	}
	return tricks, nil
}

// sorting
type ByValue []Card

func (v ByValue) Len() int { return len(v) }

func (v ByValue) Less(i, j int) bool {
	return true
}

func (v ByValue) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

func (c *Card) IsTrump() bool {
	return c.IsTrumpNumber || c.IsTrumpSuit
}
