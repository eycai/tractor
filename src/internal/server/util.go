package server

import (
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func (s *Server) generateUserID() string {
	id := randomString(s.IDLength)
	_, ok := s.Users[id]
	for ok {
		id = randomString(s.IDLength)
		_, ok = s.Users[id]
	}
	return id
}

func (s *Server) generateRoomID() string {
	id := randomString(s.IDLength)
	_, ok := s.Rooms[id]
	for ok {
		id = randomString(s.IDLength)
		_, ok = s.Rooms[id]
	}
	return id
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String() // E.g. "ExcbsVQs"
}

func getUserID(w http.ResponseWriter, r *http.Request) string {
	id, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "no user ID set in cookie", http.StatusBadRequest)
		return ""
	}
	return id.Value
}
