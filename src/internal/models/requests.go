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

type FlipCardsRequest struct {
	Card     Card `json:"card"`
	NumCards int  `json:"numCards"`
}

type SetKittyRequest struct {
	Kitty []Card `json:"kitty"`
}

type PlayCardsRequest struct {
	Cards [][]Card `json:"cards"`
}

type GameStateResponse struct {
	Room *Room `json:"room"`
	User *User `json:"user"`
}

type UpdateEvent struct {
	Room  *Room  `json:"room"`
	User  *User  `json:"user"`
	Event *Event `json:"event"`
}

type EmptyResponse struct {
}

type EndRoundEvent struct {
	Kitty       []Card `json:"kitty"`
	KittyPoints int    `json:"kittyPoints"`
	TotalPoints int    `json:"totalPoints"`
}

type Event struct {
	Name string `json:"name"`
}
