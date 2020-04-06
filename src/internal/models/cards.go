package models

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

// Card a playing card
type Card struct {
	Value         int  `json:"value"`
	Suit          Suit `json:"suit"`
	isTrumpSuit   bool
	isTrumpNumber bool
	gameValue     int
}

// Deck a standard deck of cards
type Deck struct {
	Cards []Card
}

// Suit the suit of a card
type Suit string

// Suits
const (
	Spade   Suit = "SPADE"
	Diamond Suit = "DIAMOND"
	Heart   Suit = "HEART"
	Club    Suit = "CLUB"
	Joker   Suit = "JOKER"
)

// Suits a list of suits
var Suits = []Suit{Spade, Diamond, Heart, Club, Joker}

// Pattern a type of trick
type Pattern string

// Patterns
const (
	NOfAKind Pattern = "N_OF_A_KIND"
	Tractor  Pattern = "TRACTOR"
)

// Trick a group of cards played together
type Trick struct {
	Pattern               Pattern
	LargestCard           Card
	NumCards              int
	TractorNumConsecutive int
	Suit                  Suit
	IsTrump               bool
}

// Matches returns true if the two cards are the same
func (c *Card) Matches(card Card) bool {
	return c.Value == card.Value && c.Suit == card.Suit
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

// BeatsLead checks if a given play can be beaten by a given hand.
func BeatsLead(cards [][]Card, hand []Card) (bool, error) {
	for _, t := range cards {
		trick, err := ParseTrick(t)
		if err != nil {
			return false, err
		}
		suitCards := getCardsOfSuit(trick.Suit, hand)
		for i := 0; i < len(suitCards)-trick.NumCards+1; i++ {
			c, err := ParseTrick(suitCards[i : i+trick.NumCards])
			if err != nil {
				continue
			}
			if typesMatch(c, trick) && c.LargestCard.GameValue() > trick.LargestCard.GameValue() {
				return true, nil
			}
		}
	}
	return false, nil
}

// SmallestPlay gives the smallest play in a list of tricks.
func SmallestPlay(cards [][]Card) ([][]Card, error) {
	tricks, err := GetTricks(cards)
	if err != nil {
		return [][]Card{}, err
	}

	trickMap := make(map[Trick]int)
	for i, t := range tricks {
		trickMap[t] = i
	}

	sort.Sort(ByType(tricks))
	return [][]Card{
		cards[trickMap[tricks[0]]],
	}, nil
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
				}
				trick.LargestCard = cards[i+1]
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
func GetTricks(cards [][]Card) ([]Trick, error) {
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
		if !typesMatch(t, next[i]) || t.LargestCard.gameValue >= next[i].LargestCard.gameValue {
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
// NOfAKind smaller than tractor,
// shorter length is smaller,
// shorter sublength is smaller
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

// GameValue returns the game value of the card
func (c *Card) GameValue() int {
	return c.gameValue
}

// IsTrumpNumber checks if the card is the trump number
func (c *Card) IsTrumpNumber() bool {
	return c.isTrumpNumber
}

// IsTrumpSuit checks if the card is the trump suit
func (c *Card) IsTrumpSuit() bool {
	return c.isTrumpSuit
}

// WithGameValues returns an updated card with the game values set based on vals.
func (c *Card) WithGameValues(vals map[Card]int) Card {
	c.gameValue = vals[*c]
	return *c
}

// WithTrump returns an updated card with the trump set based on n and s.
func (c *Card) WithTrump(n int, s Suit) Card {
	c.isTrumpNumber = (c.Value == n)
	c.isTrumpSuit = (c.Suit == s || c.Suit == Joker)
	return *c
}

func typesMatch(a Trick, b Trick) bool {
	return a.Pattern == b.Pattern &&
		a.NumCards == b.NumCards &&
		a.TractorNumConsecutive == b.TractorNumConsecutive
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

func numCardsOfSuit(hand []Card, suit Suit) int {
	n := 0
	for _, c := range hand {
		if c.Suit == suit {
			n++
		}
	}
	return n
}

func getCardsOfSuit(suit Suit, hand []Card) []Card {
	cards := []Card{}
	for _, c := range hand {
		if c.Suit == suit && !c.IsTrump() {
			cards = append(cards, c)
		}
	}
	sort.Sort(sort.Reverse(ByValue(cards)))
	return cards
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

func (d *Deck) shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}
