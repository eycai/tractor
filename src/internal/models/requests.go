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

type JoinRoomRequest struct {
	RoomID string `json:"roomId"`
}

type PlayerJoinedEvent struct {
	UserID string `json:"userId"`
}

type PlayerLeftEvent struct {
	UserID string `json:"userId"`
}

type LeaveRoomRequest struct {
	RoomID string `json:"roomId"`
}

type RoomCreatedEvent struct {
	RoomID string `json:"roomId"`
}
