package models

// Room a room in the tractor game.
type Room struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Users    []*UserStatus `json:"users"`
	Host     string        `json:"host"`
	Game     *Game         `json:"game"`
	Capacity int           `json:"capacity"`
}

// UserStatus metadata about the user's current status.
type UserStatus struct {
	Username  string `json:"username"`
	Connected bool   `json:"connected"`
}

// HasUser returns true if the room has a user with username user.
func (r *Room) HasUser(user string) bool {
	for _, u := range r.Users {
		if u.Username == user {
			return true
		}
	}
	return false
}

// DrawOrder returns the drawing order for the room.
func (r *Room) DrawOrder() []string {
	users := make([]string, len(r.Users))
	first := r.Game.Turn
	firstIndex := 0
	for i, u := range r.Users {
		users[i] = u.Username
		if u.Username == first {
			firstIndex = i
		}
	}
	users = append(users[firstIndex:], users[0:firstIndex]...)
	return users
}
