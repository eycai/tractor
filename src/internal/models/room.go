package models

type Room struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Users    []*UserStatus `json:"users"`
	Host     string        `json:"host"`
	Game     *Game         `json:"game"`
	Capacity int           `json:"capacity"`
}

type UserStatus struct {
	Username  string `json:"username"`
	Connected bool   `json:"connected"`
}

func (r *Room) HasUser(user string) bool {
	for _, u := range r.Users {
		if u.Username == user {
			return true
		}
	}
	return false
}
