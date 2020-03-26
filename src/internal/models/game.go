package models

import "log"

type Game struct {
	Players     []*Player         `json:"players"`
	Turn        string            `json:"turn"`
	TrumpSuit   Suit              `json:"trumpSuit"`
	TrumpNumber int               `json:"trumpNumber"`
	Banker      string            `json:"banker"`
	CardsInPlay map[string][]Card `json:"cardsInPlay"`
}

func (g *Game) getDeck() Deck {
	numDecks := len(g.Players)
	deck := make([]Card, numDecks*54)
	for _, s := range Suits {
		if s == Joker {
			newCards := make([]Card, 2*numDecks)
			for i := 0; i < numDecks; i++ {
				newCards[i] = Card{Value: 1, Suit: s}
				newCards[i+1] = Card{Value: 2, Suit: s}
			}
		} else {
			for i := 1; i <= 13; i++ {
				newCards := make([]Card, numDecks)
				for j := range newCards {
					newCards[j] = Card{Value: i, Suit: s}
				}
				deck = append(deck, newCards...)
			}
		}
	}
	d := Deck{Cards: deck}
	d.shuffle()
	log.Printf("deck: %v", d)
	return d
}
