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
