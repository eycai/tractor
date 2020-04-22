package models

import (
	"fmt"
	"log"
	"math"
	"sort"
)

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
	drawOrder            []string
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
		g.TrumpNumCardsFlipped = numCards
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
		g.cardValues = make(map[Card]int)
		g.GetCardValues()
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

	for _, p := range g.Players {
		p.ResetCards()
	}

	if g.isFirstRound() {
		bankerInd := 0
		for i, u := range g.drawOrder {
			if u == g.Banker {
				bankerInd = i
				break
			}
		}

		for i, u := range g.drawOrder {
			if i%2 == bankerInd%2 {
				g.Players[u].Team = Bosses
			} else {
				g.Players[u].Team = Peasants
			}
		}
	}
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
	if !HasCards(hand, cards) {
		return false
	}
	cards = g.GetUpdatedPlays(cards)
	hand = g.GetUpdatedCards(hand)
	firstPlay := g.GetUpdatedPlays(g.Players[g.firstInTrick].CardsPlayed)
	log.Printf("first play: %v", firstPlay)
	tricks, err := GetTricks(cards)
	if err != nil {
		cards = GetFallback(cards)
		tricks, err = GetTricks(cards)
		if err != nil {
			log.Printf("can't parse play")
			return false
		}

	}

	// is the first play
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

	// is not the first play
	return IsValidPlay(firstPlay, cards, hand)
}

// PlayCards allows the user to play cards in a trick.
// assumes that the plays are valid.
func (g *Game) PlayCards(user string, cards [][]Card, otherHands [][]Card) (TrickStatus, [][]Card, error) {
	if user != g.Turn {
		return PlayingTrick, cards, fmt.Errorf("incorrect user")
	}
	cards = g.GetUpdatedPlays(cards)

	trick, err := GetTricks(cards)
	if err != nil {
		cards = GetFallback(cards)
		print("FALLING BACK %v", cards)
		trick, err = GetTricks(cards)
		if err != nil {
			print("ERROR PARSING %v", cards)
			return PlayingTrick, cards, err
		}
	}

	if user == g.firstInTrick && len(cards) > 1 {
		// check if invalid lead, and force to play the smallest trick.
		for _, h := range otherHands {
			h = g.GetUpdatedCards(h)
			invalid, smallest, err := BeatsLead(cards, h)
			if err != nil {
				return PlayingTrick, cards, err
			}
			if invalid {
				cards = smallest
				trick, err = GetTricks(cards)
				if err != nil {
					cards = GetFallback(cards)
					print("FALLING BACK %v", cards)
					trick, err = GetTricks(cards)
					if err != nil {
						print("ERROR PARSING %v", cards)
						return PlayingTrick, cards, err
					}
				}
			}
		}
	}

	// set current winner
	if user == g.firstInTrick || NextTrickWins(g.winningTrick, trick) {
		g.SetWinningTrick(user, trick)
	}

	g.Players[user].CardsPlayed = cards
	if !g.trickEnded() {
		g.setNextTurn()
		return PlayingTrick, cards, nil
	}
	return TrickEnded, cards, nil
}

// EndState returns the ending state of the game.
func (g *Game) EndState() (int, int, Team) {
	levels, points, kittyPoints := g.peasantResults()
	team := Peasants
	if levels < 0 {
		team = Bosses
	}
	return points, kittyPoints, team
}

// SetWinningTrick sets the user and trick that is currently winning
func (g *Game) SetWinningTrick(user string, trick []Trick) {
	g.currentWinner = user
	g.winningTrick = trick
}

// EndRound ends the round and resets the game for the next round.
func (g *Game) EndRound() {
	g.round++
	peasantLevels, _, _ := g.peasantResults()
	winningTeam := Bosses
	if peasantLevels >= 0 {
		winningTeam = Peasants
	}

	// set levels
	for _, p := range g.Players {
		p.Points = 0
		if winningTeam == Peasants {
			if p.Team == Peasants {
				p.setLevel(peasantLevels)
			}
		} else {
			if p.Team == Bosses {
				p.setLevel(0 - peasantLevels)
			}
		}
		p.ResetCards()
	}

	// switch teams
	if winningTeam == Peasants {
		for _, p := range g.Players {
			p.SwitchTeam()
		}
	}

	// set banker
	bankerInd := 0
	for i := range g.drawOrder {
		if g.drawOrder[i] == g.Banker {
			bankerInd = i
		}
	}

	for i := 1; i <= len(g.drawOrder); i++ {
		ind := (i + bankerInd) % len(g.drawOrder)
		if g.Players[g.drawOrder[ind]].Team == Bosses {
			g.setBanker(g.drawOrder[ind])
			g.TrumpNumber = g.Players[g.Banker].Level
			break
		}
	}

	// reset game
	g.GamePhase = Start
	g.TrumpNumCardsFlipped = 0
	g.TrumpFlipUser = ""
	g.cardValues = make(map[Card]int)
	g.kitty = make([]Card, 0)
	g.drawOrder = make([]string, len(g.Players))
	g.winningTrick = []Trick{}
	g.currentWinner = ""
	g.firstInTrick = g.Banker
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

// GetDeck gets a plain deck.
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
	// for sorting while drawing
	deck = g.GetUpdatedCards(deck)
	d := Deck{Cards: deck}
	d.shuffle()
	return d
}

// SetTrickStarter sets the person to start the trick
func (g *Game) SetTrickStarter(user string) {
	g.firstInTrick = user
}

// SetDrawOrder sets the drawing order
func (g *Game) SetDrawOrder(order []string) {
	g.drawOrder = order
	for i, u := range order {
		g.Players[u].setOrder(i)
	}
}

// GetCurrentWinner gets the user currently winning the round
func (g *Game) GetCurrentWinner() string {
	return g.currentWinner
}

func (g *Game) setNextTurn() {
	current := g.Players[g.Turn].drawOrder
	next := (current + 1) % len(g.Players)
	g.Turn = g.drawOrder[next]
}

func (g *Game) setBanker(user string) {
	g.Banker = user
	g.firstInTrick = user
	g.Turn = user
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
	if g.Players[g.currentWinner].Team == Peasants {
		g.Players[g.currentWinner].Points += points
	}
}

func (g *Game) EndTrick() {
	g.distributePoints()
	g.Turn = g.currentWinner
	g.firstInTrick = g.currentWinner
}

func (g *Game) trickEnded() bool {
	for _, p := range g.Players {
		if len(p.CardsPlayed) == 0 {
			return false
		}
	}
	return true
}

// level, peasant points, kitty points
func (g *Game) peasantResults() (int, int, int) {
	peasantPoints := 0
	for _, p := range g.Players {
		if p.Team == Peasants {
			peasantPoints += p.Points
		}
	}

	kittyPoints := 0
	if g.Players[g.currentWinner].Team == Peasants {
		mult := 0
		for _, t := range g.winningTrick {
			mult += int(math.Exp2(float64(t.NumCards)))
		}
		kittyPoints = GetPoints([][]Card{g.kitty}) * mult
	}

	totalPoints := peasantPoints + kittyPoints
	numDecks := len(g.Players) / 2
	pointsPerLevel := numDecks * 20
	if totalPoints == 0 {
		return -3, peasantPoints, kittyPoints
	}
	return (totalPoints / pointsPerLevel) - 2, peasantPoints, kittyPoints
}
