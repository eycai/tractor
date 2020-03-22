package models

type Room struct {
	ID     string   `json:"id"`
	Users  []string `json:"users"`
	HostID string   `json:"hostId"`
	Game   Game     `json:"Game"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	SocketID string `json:"socketId"`
	RoomID   string `json:"roomId"`
}
