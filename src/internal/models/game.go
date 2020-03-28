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
	// suitValue := 0
	for _, s := range Suits {
		if s == Joker {
			newCards := make([]Card, 2*numDecks)
			for i := 0; i < numDecks; i++ {
				newCards[i] = Card{Value: 1, Suit: s, gameValue: 53, IsTrumpSuit: true}
				newCards[i+1] = Card{Value: 2, Suit: s, gameValue: 54, IsTrumpSuit: true}
			}
		} else {
			for i := 1; i <= 13; i++ {
				newCards := make([]Card, numDecks)
				// gameValue := k*13 + i
				// if i == g.TrumpNumber && s == g.TrumpSuit {
				// 	gameValue = 52
				// } else if i == g.TrumpNumber {
				// 	gameValue = 51
				// }
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
