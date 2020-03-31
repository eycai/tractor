package models

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

type Card struct {
	Value int  `json:"value"`
	Suit  Suit `json:"suit"`
	// IsTrumpSuit   bool `json:"isTrumpSuit"`
	// IsTrumpNumber bool `json:"isTrumpNumber"`
	// GameValue     int  `json:"gameValue"`
	isTrumpSuit   bool
	isTrumpNumber bool
	gameValue     int
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

func numCardsOfSuit(hand []Card, suit Suit) int {
	n := 0
	for _, c := range hand {
		if c.Suit == suit {
			n++
		}
	}
	return n
}

func numTrumpCards(hand []Card) int {
	n := 0
	for _, c := range hand {
		if c.IsTrump() {
			n++
		}
	}
	return n
}

func cardList(hand [][]Card) []Card {
	cards := []Card{}
	for _, c := range hand {
		cards = append(cards, c...)
	}
	return cards
}

func lenPlay(cards [][]Card) int {
	n := 0
	for _, c := range cards {
		n += len(c)
	}
	return n
}

// GetPoints returns the number of points in a hand.
func GetPoints(hand [][]Card) int {
	points := 0
	for _, trick := range hand {
		for _, c := range trick {
			if c.Value == 5 {
				points += 5
			} else if c.Value == 10 || c.Value == 13 {
				points += 10
			}
		}
	}
	return points
}

// IsValidPlay returns true if the two plays match in length, and suit is valid
func IsValidPlay(prev [][]Card, next [][]Card, hand []Card) bool {
	if lenPlay(prev) != lenPlay(next) {
		return false
	}

	// used fewer cards of suit than min(available, length of hand)
	if prev[0][0].IsTrump() && numTrumpCards(cardList(next)) <
		int(math.Min(float64(numTrumpCards(hand)), float64(len(hand)))) {
		return false
	}
	if !prev[0][0].IsTrump() && numCardsOfSuit(cardList(next), prev[0][0].Suit) <
		int(math.Min(float64(numCardsOfSuit(hand, prev[0][0].Suit)), float64(len(hand)))) {
		return false
	}

	return true
}

// ParseTrick parses a list of cards into a trick.
func ParseTrick(cards []Card) (Trick, error) {
	sort.Sort(sort.Reverse(ByValue(cards)))
	trick := Trick{
		Pattern:     NOfAKind,
		NumCards:    len(cards),
		Suit:        cards[0].Suit,
		IsTrump:     cards[0].IsTrump(),
		LargestCard: cards[0],
	}

	numConsecutive := 1
	for i := 0; i < len(cards)-1; i++ {
		if !(trick.IsTrump == cards[i+1].IsTrump()) {
			return trick, fmt.Errorf("either all cards should be trump, or none")
		}
		if !(trick.IsTrump && cards[i+1].IsTrump()) && trick.Suit != cards[i+1].Suit {
			return trick, fmt.Errorf("all suits should be the same")
		}
		if cards[i] != cards[i+1] {
			if trick.TractorNumConsecutive != 0 && trick.TractorNumConsecutive != numConsecutive {
				return trick, fmt.Errorf("tractor incorrect length")
			} else if !IsConsecutive(cards[i], cards[i+1]) {
				if !(cards[0].Value == 1 && cards[len(cards)-1].Value == 2) {
					return trick, fmt.Errorf("tractor not consecutive")
				} else {
					trick.LargestCard = cards[i+1]
				}
			}
			trick.Pattern = Tractor
			trick.TractorNumConsecutive = numConsecutive
			numConsecutive = 1
		} else {
			numConsecutive++
		}
	}
	if trick.TractorNumConsecutive != 0 && numConsecutive != trick.TractorNumConsecutive {
		return trick, fmt.Errorf("tractor incorrect length")
	}
	return trick, nil
}

// IsConsecutive determines if two cards are consecutive. It assumes that the cards are
// of the same suit.
func IsConsecutive(a Card, b Card) bool {
	return math.Abs(float64(a.gameValue-b.gameValue)) == 1
}

// GetTricks parses tricks out of a play.
func GetTricks(cards [][]Card, trumpSuit Suit, trumpNumber int) ([]Trick, error) {
	tricks := make([]Trick, len(cards))
	for i, t := range cards {
		trick, err := ParseTrick(t)
		if err != nil {
			return tricks, err
		}
		tricks[i] = trick
	}
	return tricks, nil
}

func trickSuitsMatch(t []Trick) bool {
	suitToMatch := t[0].Suit
	isTrump := t[0].IsTrump
	for _, s := range t {
		if isTrump {
			if !s.IsTrump {
				return false
			}
		} else if s.Suit != suitToMatch {
			return false
		}
	}
	return true
}

func typesMatch(a Trick, b Trick) bool {
	return a.Pattern == b.Pattern &&
		a.NumCards == b.NumCards &&
		a.TractorNumConsecutive == b.TractorNumConsecutive
}

// NextTrickWins returns true if the second play is larger; else, false.
func NextTrickWins(prev []Trick, next []Trick) bool {
	sort.Sort(sort.Reverse(ByType(prev)))
	sort.Sort(sort.Reverse(ByType(next)))

	if len(prev) != len(next) {
		// tricks don't match
		return false
	}

	if !trickSuitsMatch(next) {
		return false
	}

	if prev[0].IsTrump && !next[0].IsTrump {
		return false
	}
	if !prev[0].IsTrump && !next[0].IsTrump && prev[0].Suit != next[0].Suit {
		// suits don't match
		return false
	}

	// only remaining: neither are trump and suits match, or
	// both are trump, or
	// b is trump but not a.
	// in all cases, check that pattern maps, and that game value is larger.

	for i, t := range prev {
		if !typesMatch(t, next[i]) || t.LargestCard.gameValue > next[i].LargestCard.gameValue {
			return false
		}
	}
	return true
}

// ByValue allows for sorting cards by game value
type ByValue []Card

func (v ByValue) Len() int           { return len(v) }
func (v ByValue) Less(i, j int) bool { return v[i].gameValue < v[j].gameValue }
func (v ByValue) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }

// ByType allows for sorting tricks by type.
// Our desired sort patterns are:
// By pattern, then by length, then by sublength, then by game value.
type ByType []Trick

func (v ByType) Len() int { return len(v) }
func (v ByType) Less(i, j int) bool {
	if v[i].Pattern != v[j].Pattern {
		return v[i].Pattern == NOfAKind
	}
	if v[i].NumCards != v[j].NumCards {
		return v[i].NumCards < v[j].NumCards
	}
	if v[i].Pattern == Tractor && v[i].TractorNumConsecutive != v[j].TractorNumConsecutive {
		return v[i].TractorNumConsecutive < v[j].TractorNumConsecutive
	}
	return v[i].LargestCard.gameValue < v[j].LargestCard.gameValue
}
func (v ByType) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

// IsTrump returns true if the card is trump
func (c *Card) IsTrump() bool {
	return c.isTrumpNumber || c.isTrumpSuit
}

func (c *Card) GameValue() int {
	return c.gameValue
}

func (c *Card) IsTrumpNumber() bool {
	return c.isTrumpNumber
}

func (c *Card) IsTrumpSuit() bool {
	return c.isTrumpSuit
}

func (c *Card) WithGameValues(vals map[Card]int) Card {
	c.gameValue = vals[*c]
	return *c
}

func (c *Card) WithTrump(n int, s Suit) Card {
	if c.Value == n {
		c.isTrumpNumber = true
	}

	if c.Suit == s || c.Suit == Joker {
		c.isTrumpSuit = true
	}
	return *c
}
