package models

type Room struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Users    []string `json:"users"`
	Host     string   `json:"host"`
	Game     *Game    `json:"game"`
	Capacity int      `json:"capacity"`
}

func (r *Room) HasUser(user string) bool {
	for _, u := range r.Users {
		if u == user {
			return true
		}
	}
	return false
}
