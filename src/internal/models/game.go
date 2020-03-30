package models

type Game struct {
	Players     map[string]*Player `json:"players"`
	Turn        string             `json:"turn"`
	TrumpSuit   Suit               `json:"trumpSuit"`
	TrumpNumber int                `json:"trumpNumber"`
	Banker      string             `json:"banker"`
	CardsInPlay map[string][]Card  `json:"cardsInPlay"`
	GamePhase   Phase              `json:"gamePhase"`
}

type Phase string

const (
	Drawing  Phase = "DRAWING"
	Playing  Phase = "PLAYING"
	EndRound Phase = "END_ROUND"
)

func (g *Game) GetDeck() Deck {
	numDecks := len(g.Players) / 2
	deck := []Card{}
	suitNum := 0
	for _, s := range Suits {
		if s == Joker {
			newCards := make([]Card, 2*numDecks)
			for i := 0; i < numDecks; i++ {
				newCards[2*i] = Card{Value: 1, Suit: s, GameValue: 50, IsTrumpSuit: true}
				newCards[2*i+1] = Card{Value: 2, Suit: s, GameValue: 51, IsTrumpSuit: true}
			}
			deck = append(deck, newCards...)
		} else {
			currentValue := suitNum * 12
			if s == g.TrumpSuit {
				currentValue = 3 * 12
			} else {
				suitNum++
			}
			for i := 1; i <= 13; i++ {
				newCards := make([]Card, numDecks)
				gameValue := currentValue
				if i == g.TrumpNumber && s == g.TrumpSuit {
					gameValue = 49
				} else if i == g.TrumpNumber {
					gameValue = 48
				} else if i == 1 {
					gameValue = currentValue + 11
				} else {
					currentValue++
				}
				for j := range newCards {
					newCards[j] = Card{Value: i, Suit: s, GameValue: gameValue}
					if s == g.TrumpSuit {
						newCards[j].IsTrumpSuit = true
					}
					if i == g.TrumpNumber {
						newCards[j].IsTrumpNumber = true
					}
				}
				deck = append(deck, newCards...)
			}
		}
	}
	d := Deck{Cards: deck}
	d.shuffle()
	return d
}
