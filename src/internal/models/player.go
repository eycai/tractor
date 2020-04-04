package models

type Player struct {
	Username    string   `json:"username"`
	Team        Team     `json:"team"`
	Level       int      `json:"level"`
	Points      int      `json:"points"`
	CardsPlayed [][]Card `json:"cardsPlayed"`
}

type Team string

const (
	Bosses   Team = "BOSSES"
	Peasants Team = "PEASANTS"
)

func (p *Player) PlayCards(cards [][]Card) {
	p.CardsPlayed = cards
}

func (p *Player) ResetCards() {
	p.CardsPlayed = [][]Card{}
}

func (p *Player) SwitchTeam() {
	if p.Team == Bosses {
		p.Team = Peasants
	} else {
		p.Team = Bosses
	}
}
