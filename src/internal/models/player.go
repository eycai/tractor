package models

// Player a player in a game of tractor
type Player struct {
	Username    string   `json:"username"`
	Team        Team     `json:"team"`
	Level       int      `json:"level"`
	Points      int      `json:"points"`
	CardsPlayed [][]Card `json:"cardsPlayed"`
	drawOrder   int
}

// Team the teams in tractor
type Team string

// Teams
const (
	Bosses   Team = "BOSSES"
	Peasants Team = "PEASANTS"
)

// PlayCards has the user play the given cards
func (p *Player) PlayCards(cards [][]Card) {
	p.CardsPlayed = cards
}

// ResetCards resets the played cards to none
func (p *Player) ResetCards() {
	p.CardsPlayed = [][]Card{}
}

// SwitchTeam sets the player's current team to the opposite team
func (p *Player) SwitchTeam() {
	if p.Team == Bosses {
		p.Team = Peasants
	} else {
		p.Team = Bosses
	}
}

func (p *Player) setOrder(i int) {
	p.drawOrder = i
}
