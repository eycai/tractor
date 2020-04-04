package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	SocketID string `json:"socketId"`
	RoomID   string `json:"roomId"`
	Hand     []Card `json:"hand"`
	Kitty    []Card `json:"kitty"`
}

func (u *User) Reset() {
	u.SocketID = ""
	u.RoomID = ""
	u.Hand = []Card{}
	u.Kitty = []Card{}
}

func (u *User) DealCard(c Card) {
	u.Hand = append(u.Hand, c)
}

func (u *User) PlayCards(cards [][]Card) {
	for _, t := range cards {
		for _, c := range t {
			u.removeCardFromHand(c)
		}
	}
}

func (u *User) removeCardFromHand(card Card) {
	for i, c := range u.Hand {
		if c.Matches(card) {
			// remove first match
			u.Hand = append(u.Hand[0:i], u.Hand[i+1:]...)
			return
		}
	}
}
