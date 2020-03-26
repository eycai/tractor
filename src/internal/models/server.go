package models

type Room struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Users    []string `json:"users"`
	Host     string   `json:"host"`
	Game     *Game    `json:"game"`
	Capacity int      `json:"capacity"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	SocketID string `json:"socketId"`
	RoomID   string `json:"roomId"`
	Hand     []Card `json:"hand"`
	Kitty    []Card `json:"kitty"`
}
