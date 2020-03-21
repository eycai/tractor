package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/eycai/tractor/src/internal/models"
)

func (s *Server) ConnectUser(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	req := models.ConnectRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "error decoding request", http.StatusBadRequest)
		return
	}
	userID := getUserID(w, r)
	if userID == "" {
		return
	}
	s.SocketIDs[userID] = req.SocketID
	log.Printf("current sockets: %v", s.SocketIDs)

	w.WriteHeader(http.StatusOK)
}

func (s *Server) RegisterUser(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	req := models.RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "error decoding request", http.StatusBadRequest)
		return
	}

	if mapContainsKey(s.SocketIDs, req.Username) {
		http.Error(w, "username already taken", http.StatusConflict)
		return
	}

	expiration := time.Now().Add(365 * 24 * time.Hour)
	userID := s.generateUserID()
	s.UserIDs[req.Username] = userID
	s.SocketIDs[userID] = ""

	log.Printf("current users: %v", s.UserIDs)

	cookie := http.Cookie{Name: "user_id", Value: userID, Expires: expiration}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}
