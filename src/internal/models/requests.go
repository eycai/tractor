package models

type ConnectEvent struct {
	SocketID string `json:"socketId"`
}

type RegisterRequest struct {
	Username string `json:"username"`
}

type CreateRoomRequest struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

type ConnectRequest struct {
	SocketID string `json:"socketId"`
}

type JoinRoomRequest struct {
	RoomID string `json:"roomId"`
}

type LeaveRoomRequest struct {
	RoomID string `json:"roomId"`
}

type UpdateEvent struct {
	Room *Room `json:"room"`
	User *User `json:"user"`
}

type EmptyResponse struct {
}
