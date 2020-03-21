package models

type ConnectEvent struct {
	SocketID string `json:"socketId"`
}

type RegisterRequest struct {
	Username string `json:"username"`
}

type ConnectRequest struct {
	SocketID string `json:"socketId"`
}
