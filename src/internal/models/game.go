package models

import "sort"

// Game a game of tractor
type Game struct {
	Players              map[string]*Player `json:"players"`
	Turn                 string             `json:"turn"`
	TrumpSuit            Suit               `json:"trumpSuit"`
	TrumpNumber          int                `json:"trumpNumber"`
	Banker               string             `json:"banker"`
	GamePhase            Phase              `json:"gamePhase"`
	TrumpFlipUser        string             `json:"trumpFlipUser"`
	TrumpNumCardsFlipped int                `json:"trumpNumCardsFlipped"`
	round                int
	kitty                []Card
	firstInTrick         string
	currentWinner        string
	winningTrick         []Trick
	cardValues           map[Card]int
}

// Phase the phase of the game
type Phase string

// Phases
const (
	Start           Phase = "START"
	Drawing         Phase = "DRAWING"
	DrawingComplete Phase = "DRAWING_COMPLETE"
	SetKitty        Phase = "SET_KITTY"
	Playing         Phase = "PLAYING"
	EndRound        Phase = "END_ROUND"
)

// TrickStatus the status of the current trick
type TrickStatus string

// TrickStatuses
const (
	PlayingTrick TrickStatus = "PLAYING_TRICK"
	TrickEnded   TrickStatus = "TRICK_ENDED"
)

// FlipCard allows users to flip a card to declare trump
func (g *Game) FlipCard(c Card, numCards int, user string) bool {
	if g.GamePhase != Drawing && g.GamePhase != DrawingComplete {
		// can't set after done drawing
		return false
	}

	if user == g.TrumpFlipUser && c.Suit != g.TrumpSuit {
		return false
	}

	if c.Value != g.TrumpNumber && c.Suit != Joker {
		return false
	}

	if user == g.TrumpFlipUser {
		// reinforce
		g.TrumpNumCardsFlipped += numCards
		if g.isFirstRound() {
			g.setBanker(user)
		}
		return true
	}

	if (numCards > g.TrumpNumCardsFlipped && c.Value == g.TrumpNumber) ||
		numCards >= g.TrumpNumCardsFlipped && c.Suit == Joker {
		// overflip
		g.TrumpSuit = c.Suit
		g.TrumpNumCardsFlipped = numCards
		g.TrumpFlipUser = user
		if g.isFirstRound() {
			g.setBanker(user)
		}
		return true
	}

	return false
}

// KittySize gets the kitty size for the game
func (g *Game) KittySize() int {
	numPlayers := len(g.Players)
	numCards := (len(g.Players) / 2) * 54
	kitty := numCards % numPlayers
	for kitty <= 4 {
		kitty += numPlayers
	}
	return kitty
}

// SetKitty sets the kitty
func (g *Game) SetKitty(k []Card) {
	g.kitty = k
}

// GetKitty gets the kitty. This should only be used to get the kitty for the banker.
func (g *Game) GetKitty() []Card {
	k := g.kitty
	g.kitty = []Card{}
	return k
}

// GetUpdatedCards updates the cards with game values and trump
func (g *Game) GetUpdatedCards(cards []Card) []Card {
	vals := g.GetCardValues()
	newCards := make([]Card, len(cards))
	for i, c := range cards {
		newCards[i] = c.WithGameValues(vals)
		newCards[i] = newCards[i].WithTrump(g.TrumpNumber, g.TrumpSuit)
	}
	sort.Sort(ByValue(newCards))
	return newCards
}

// GetUpdatedPlays updates the cards with the game values and trump
func (g *Game) GetUpdatedPlays(cards [][]Card) [][]Card {
	for i, c := range cards {
		cards[i] = g.GetUpdatedCards(c)
	}
	return cards
}

// IsValidPlayForGame determines if the play is valid given current game circumstances
func (g *Game) IsValidPlayForGame(cards [][]Card, hand []Card) bool {
	cards = g.GetUpdatedPlays(cards)
	firstPlay := g.GetUpdatedPlays(g.Players[g.firstInTrick].CardsPlayed)
	tricks, err := GetTricks(cards)
	if err != nil {
		return false
	}

	if len(firstPlay) == 0 {
		if len(cards) == 1 {
			return true
		}

		suit := tricks[0].Suit
		for _, t := range tricks {
			if t.IsTrump {
				return false
			}
			if t.Suit != suit {
				return false
			}
		}
		return true
	}

	return IsValidPlay(firstPlay, cards, hand)
}

func (g *Game) PlayCards(user string, cards [][]Card, otherHands [][]Card) (TrickStatus, error) {
	cards = g.GetUpdatedPlays(cards)
	trick, err := GetTricks(cards)
	if err != nil {
		return PlayingTrick, err
	}

	if user == g.firstInTrick || NextTrickWins(g.winningTrick, trick) {
		g.currentWinner = user
		g.winningTrick = trick
	}

	if user == g.firstInTrick && len(cards) > 1 {
		for _, h := range otherHands {
			invalid, err := BeatsLead(cards, h)
			if err != nil {
				return PlayingTrick, err
			}

			if invalid {
				cards, err = SmallestPlay(cards)
				if err != nil {
					return PlayingTrick, err
				}
			}
		}
	}
	g.Players[user].CardsPlayed = cards
	if !g.trickEnded() {
		return PlayingTrick, nil
	}
	g.endTrick()
	return TrickEnded, nil
}

func (g *Game) EndRound(playOrder []string) {
	g.round++
	peasantLevels := g.peasantLevels()
	winningTeam := Bosses
	if peasantLevels >= 0 {
		winningTeam = Peasants
	}

	// set levels
	for _, p := range g.Players {
		p.Points = 0
		if winningTeam == Peasants {
			if p.Team == Peasants {
				p.Level += peasantLevels
			}
		} else {
			if p.Team == Bosses {
				p.Level -= peasantLevels
			}
		}
	}

	// switch teams
	if winningTeam == Peasants {
		for _, p := range g.Players {
			p.SwitchTeam()
		}
	}

	// set banker
	bankerInd := 0
	for i := range playOrder {
		if playOrder[i] == g.Banker {
			bankerInd = i
		}
	}

	for i := 1; i <= len(playOrder); i++ {
		ind := (i + bankerInd) % len(playOrder)
		if g.Players[playOrder[ind]].Team == Bosses {
			g.Banker = playOrder[ind]
		}
	}
}

// GetCardValues gets a map from plain card to card with trump updates after trump is set.
func (g *Game) GetCardValues() map[Card]int {
	if len(g.cardValues) != 0 {
		return g.cardValues
	}
	cardValues := make(map[Card]int)
	suitNum := 0
	for _, s := range Suits {
		if s == Joker {
			cardValues[Card{Value: 1, Suit: s}] = 50
			cardValues[Card{Value: 2, Suit: s}] = 51
		} else {
			currentValue := suitNum * 12
			if s == g.TrumpSuit {
				currentValue = 3 * 12
			} else {
				suitNum++
			}
			for i := 1; i <= 13; i++ {
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
				cardValues[Card{Value: i, Suit: s}] = gameValue
			}
		}
	}
	g.cardValues = cardValues
	return cardValues
}

// GetDeck gets a plain deck, with no trump or game value fields set.
func (g *Game) GetDeck() Deck {
	numDecks := len(g.Players) / 2
	deck := []Card{}
	for _, s := range Suits {
		if s == Joker {
			newCards := make([]Card, 2*numDecks)
			for i := 0; i < numDecks; i++ {
				newCards[2*i] = Card{Value: 1, Suit: s}
				newCards[2*i+1] = Card{Value: 2, Suit: s}
			}
			deck = append(deck, newCards...)
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
	return d
}

func (g *Game) isFirstRound() bool {
	return g.round == 0
}

func (g *Game) distributePoints() {
	points := 0
	for _, p := range g.Players {
		points += GetPoints(p.CardsPlayed)
		p.ResetCards()
	}
	g.Players[g.currentWinner].Points += points
}

// SetTrickStarter sets the person to start the trick
func (g *Game) SetTrickStarter(user string) {
	g.firstInTrick = user
}

func (g *Game) setBanker(user string) {
	g.Banker = user
	g.firstInTrick = user
}

func (g *Game) endTrick() {
	g.distributePoints()
	g.Turn = g.currentWinner
	g.firstInTrick = g.currentWinner
	g.currentWinner = ""
	g.winningTrick = []Trick{}
}

func (g *Game) trickEnded() bool {
	for _, p := range g.Players {
		if len(p.CardsPlayed) == 0 {
			return true
		}
	}
	return false
}

func (g *Game) peasantLevels() int {
	peasantPoints := 0
	for _, p := range g.Players {
		if p.Team == Peasants {
			peasantPoints += p.Points
		}
	}
	numDecks := len(g.Players) / 2
	pointsPerLevel := numDecks * 20
	return (peasantPoints / pointsPerLevel) - 2
}
