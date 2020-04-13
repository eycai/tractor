package models

// User a user in the tractor game
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	SocketID string `json:"socketId"`
	RoomID   string `json:"roomId"`
	Hand     []Card `json:"hand"`
	Kitty    []Card `json:"kitty"`
}

// Reset resets the user completely
func (u *User) Reset() {
	u.SocketID = ""
	u.RoomID = ""
	u.Hand = []Card{}
	u.Kitty = []Card{}
}

// DealCard deals card c into the user's hands
func (u *User) DealCard(c Card, g *Game) {
	u.Hand = append(u.Hand, c)
	u.Hand = g.GetUpdatedCards(u.Hand)
}

// PlayCards removes a play's cards from the user's hand.
func (u *User) PlayCards(cards [][]Card) {
	for _, t := range cards {
		for _, c := range t {
			u.removeCardFromHand(c)
		}
	}
}

// UpdateWithKitty lets the user set a new kitty.
func (u *User) UpdateWithKitty(newKitty []Card) {
	for _, c := range newKitty {
		u.removeCardFromHand(c)
	}
	u.Kitty = newKitty
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
