package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/eycai/tractor/src/internal/models"
)

func (s *Server) JoinRoom(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.mu.Unlock()

	req := models.JoinRoomRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "error decoding request", http.StatusBadRequest)
		return
	}

	userID := getUserID(w, r)
	if userID == "" {
		return
	}

	s.Users[userID].RoomID = req.RoomID
	s.Rooms[req.RoomID].Users = append(s.Rooms[req.RoomID].Users, userID)
	s.broadcastToRoom(req.RoomID, "player_joined", models.PlayerJoinedEvent{
		UserID: userID,
	})
}

func (s *Server) LeaveRoom(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	req := models.LeaveRoomRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "error decoding request", http.StatusBadRequest)
		return
	}
	userID := getUserID(w, r)
	if userID == "" {
		return
	}

	s.WSServer.leaveRoom(s.Users[userID].SocketID, req.RoomID)
	s.broadcastToRoom(req.RoomID, "player_left", models.PlayerLeftEvent{
		UserID: userID,
	})
}

func (s *Server) GetRooms(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rooms := make([]models.Room, 4)
	rooms[0] = models.Room{
		ID:     fmt.Sprintf("%d", 0),
		Users:  []string{"a", "b", "c"},
		HostID: "a",
	}
	rooms[1] = models.Room{
		ID:     fmt.Sprintf("%d", 1),
		Users:  []string{"e", "f"},
		HostID: "f",
	}
	rooms[2] = models.Room{
		ID:     fmt.Sprintf("%d", 2),
		Users:  []string{"g"},
		HostID: "g",
	}
	rooms[3] = models.Room{
		ID:     fmt.Sprintf("%d", 3),
		Users:  []string{"h", "i", "j", "k"},
		HostID: "j",
		Game: models.Game{
			Players:     []models.Player{},
			Kitty:       []models.Card{},
			TurnUserID:  "j",
			TrumpSuit:   models.Diamond,
			TrumpNumber: 2,
			BankerID:    "j",
			CardsInPlay: make(map[string][]models.Card),
		},
	}

	roomsJSON, err := json.Marshal(&rooms)
	if err != nil {
		http.Error(w, "error marshalling rooms", http.StatusInternalServerError)
	}
	w.Write(roomsJSON)
}

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
	s.Users[userID].SocketID = req.SocketID
	log.Printf("current users: %v", s.Users)
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

	if _, ok := s.UserIDs[req.Username]; ok {
		http.Error(w, "username already taken", http.StatusConflict)
		return
	}

	expiration := time.Now().Add(365 * 24 * time.Hour)
	userID := s.generateUserID()
	s.UserIDs[req.Username] = userID
	s.Users[userID] = &models.User{
		ID:       userID,
		Username: req.Username,
	}

	log.Printf("current users: %v", s.UserIDs)

	cookie := http.Cookie{Name: "user_id", Value: userID, Expires: expiration}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) CreateRoom(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(w, r)
	if userID == "" {
		return
	}
	roomID := s.generateRoomID()
	s.WSServer.emit(s.Users[userID].SocketID, "room_created", models.RoomCreatedEvent{
		RoomID: roomID,
	})
	s.Rooms[roomID] = &models.Room{
		ID:     roomID,
		Users:  []string{userID},
		HostID: userID,
	}
	s.Users[userID].RoomID = roomID
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(w, r)
	if userID == "" {
		return
	}

	user := s.Users[userID]
	userJSON, err := json.Marshal(&user)
	if err != nil {
		http.Error(w, "error marshalling user", http.StatusInternalServerError)
	}
	w.Write(userJSON)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) StartGame(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(w, r)
	if userID == "" {
		return
	}

	roomID := s.Users[userID].RoomID
	s.broadcastToRoom(roomID, "game_started", "started")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) broadcastToRoom(roomID string, eventName string, event interface{}) {
	for _, id := range s.Rooms[roomID].Users {
		s.WSServer.emit(s.Users[id].SocketID, eventName, event)
	}
}
