package models

type Player struct {
	Username string `json:"username"`
	Team     Team   `json:"team"`
	Level    int    `json:"level"`
	Points   int    `json:"points"`
}

type Team string

const (
	Bosses   Team = "BOSSES"
	Peasants Team = "PEASANTS"
)
