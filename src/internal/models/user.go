package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	SocketID string `json:"socketId"`
	RoomID   string `json:"roomId"`
	Hand     []Card `json:"hand"`
	Kitty    []Card `json:"kitty"`
}
