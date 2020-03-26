package server

import (
	"math/rand"
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

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func indexOf(slice []string, s string) int {
	for i := range slice {
		if slice[i] == s {
			return i
		}
	}
	return -1
}
