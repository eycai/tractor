package server

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/eycai/tractor/src/internal/models"
)

func setCookie(w http.ResponseWriter, name string, value string) {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: name, Value: value, Expires: expiration}
	http.SetCookie(w, &cookie)
}

func removeCookie(w http.ResponseWriter, name string, value string) {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: name, Value: value, Expires: expiration}
	log.Printf("removed cookie")
	http.SetCookie(w, &cookie)
}

func (s *Server) processPostRequest(w http.ResponseWriter, r *http.Request, req interface{}) (string, error) {
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, "error decoding request", http.StatusBadRequest)
		return "", err
	}

	userID := s.getUserID(w, r)
	if userID == "" {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return "", fmt.Errorf("invalid user id %s", userID)
	}

	return userID, nil
}

func (s *Server) broadcastUpdate(roomID string, eventName string) {
	if _, ok := s.Rooms[roomID]; !ok {
		return
	}

	for _, user := range s.Rooms[roomID].Users {
		userID := s.UserIDs[user.Username]
		s.emitUpdateToUser(userID, eventName)
	}
}

func (s *Server) emitUpdateToUser(userID string, updateEventName string) {
	update := models.UpdateEvent{
		User: s.Users[userID],
		Event: &models.Event{
			Name: updateEventName,
		},
	}

	// update room if it exists
	if s.Users[userID].RoomID != "" {
		if room, ok := s.Rooms[s.Users[userID].RoomID]; ok {
			update.Room = room
		}
	}

	s.emitWSToUser(userID, "update", update)
}

func (s *Server) emitWSToUser(userID string, eventName string, event interface{}) {
	s.Emit(s.Users[userID].SocketID, eventName, event)
}

func returnSuccess(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	resp, err := json.Marshal(&models.EmptyResponse{})
	if err != nil {
		http.Error(w, "error writing response", http.StatusInternalServerError)
	}
	w.Write(resp)
}

func (s *Server) getUserID(w http.ResponseWriter, r *http.Request) string {
	id, err := r.Cookie("user_id")
	if err != nil {
		return ""
	}
	if _, ok := s.Users[id.Value]; !ok {
		removeCookie(w, "user_id", "")
		return ""
	}
	return id.Value
}

func (s *Server) generateUserID() string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	id := randomString(s.UserIDLength, chars)
	_, ok := s.Users[id]
	for ok {
		id = randomString(s.UserIDLength, chars)
		_, ok = s.Users[id]
	}
	return id
}

func (s *Server) generateRoomID() string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	id := randomString(s.RoomIDLength, chars)
	_, ok := s.Rooms[id]
	for ok {
		id = randomString(s.RoomIDLength, chars)
		_, ok = s.Rooms[id]
	}
	return id
}

func randomString(length int, chars []rune) string {
	rand.Seed(time.Now().UnixNano())
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String() // E.g. "ExcbsVQs"
}

func remove(slice []*models.UserStatus, s int) []*models.UserStatus {
	return append(slice[:s], slice[s+1:]...)
}

func indexOf(slice []*models.UserStatus, s string) int {
	for i := range slice {
		if slice[i].Username == s {
			return i
		}
	}
	return -1
}
