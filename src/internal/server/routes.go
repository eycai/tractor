package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/eycai/tractor/src/internal/models"
)

func (s *Server) JoinRoom(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	req := models.JoinRoomRequest{}
	userID, err := s.processPostRequest(w, r, &req)
	if err != nil {
		return
	}

	room, ok := s.Rooms[req.RoomID]
	if !ok {
		http.Error(w, "no room with given id", http.StatusBadRequest)
		return
	}

	// check capacity
	if len(room.Users) >= room.Capacity {
		http.Error(w, "room at capacity", http.StatusBadRequest)
		return
	}

	if room.Game != nil {
		http.Error(w, "game in progress", http.StatusBadRequest)
		return
	}

	s.WSServer.AddToRoom(s.Users[userID].SocketID, req.RoomID)
	s.addToRoom(userID, req.RoomID)
	s.broadcastUpdate(req.RoomID, "player_joined")
	returnSuccess(w)
}

func (s *Server) LeaveRoom(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	req := models.LeaveRoomRequest{}
	userID, err := s.processPostRequest(w, r, &req)
	if err != nil {
		return
	}

	room, ok := s.Rooms[req.RoomID]
	if !ok {
		http.Error(w, "no room with given id", http.StatusBadRequest)
		return
	}

	if room.Game != nil {
		http.Error(w, "game in progress", http.StatusBadRequest)
		return
	}

	s.WSServer.LeaveRoom(s.Users[userID].SocketID, req.RoomID)
	s.removeFromRoom(userID, req.RoomID)
	s.broadcastUpdate(req.RoomID, "player_left")
	returnSuccess(w)
}

func (s *Server) GetRooms(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	roomsList := make([]*models.Room, len(s.Rooms))
	i := 0
	for _, r := range s.Rooms {
		roomsList[i] = r
		i++
	}

	roomsJSON, err := json.Marshal(&roomsList)
	if err != nil {
		http.Error(w, "error marshalling rooms", http.StatusInternalServerError)
	}
	w.Write(roomsJSON)
}

func (s *Server) RoomInfo(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	roomID := r.URL.Query().Get("userId")
	userID := s.getUserID(w, r)
	if userID == "" {
		return
	}

	room, err := json.Marshal(s.Rooms[roomID])
	if err != nil {
		http.Error(w, "error marshalling rooms", http.StatusInternalServerError)
	}

	w.Write(room)
}

func (s *Server) RegisterUser(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// decode request
	req := models.RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "error decoding request", http.StatusBadRequest)
		return
	}

	// no empty username
	if len(req.Username) == 0 {
		http.Error(w, "invalid username", http.StatusBadRequest)
		return
	}

	// username already used
	if _, ok := s.UserIDs[req.Username]; ok {
		http.Error(w, "username already taken", http.StatusConflict)
		return
	}

	// create user id
	userID := s.generateUserID()
	s.UserIDs[req.Username] = userID
	setCookie(w, "user_id", userID)

	// add user
	s.Users[userID] = &models.User{
		ID:       userID,
		Username: req.Username,
	}

	log.Printf("current users: %v", s.UserIDs)
	returnSuccess(w)
}

func (s *Server) ConnectUser(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	req := models.ConnectRequest{}
	userID, err := s.processPostRequest(w, r, &req)
	if err != nil {
		return
	}

	// add socket map
	s.Users[userID].SocketID = req.SocketID
	log.Printf("updated the socket id to be %v", s.Users[userID].SocketID)
	log.Printf("current users: %v", s.Users)
	returnSuccess(w)
}

func (s *Server) CreateRoom(w http.ResponseWriter, r *http.Request) {
	req := models.CreateRoomRequest{}
	userID, err := s.processPostRequest(w, r, &req)
	if err != nil {
		return
	}

	room := s.createRoom(userID, req.Name, req.Capacity)
	s.WSServer.AddToRoom(s.Users[userID].SocketID, room.ID)
	roomJSON, err := json.Marshal(&room)
	if err != nil {
		http.Error(w, "error marshalling room", http.StatusInternalServerError)
	}
	w.Write(roomJSON)
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := s.getUserID(w, r)
	user := &models.User{}
	if userID != "" {
		user = s.Users[userID]
	}

	userJSON, err := json.Marshal(&user)
	if err != nil {
		http.Error(w, "error marshalling user", http.StatusInternalServerError)
	}
	w.Write(userJSON)
}

func (s *Server) StartGame(w http.ResponseWriter, r *http.Request) {
	userID := s.getUserID(w, r)
	if userID == "" {
		return
	}

	roomID := s.Users[userID].RoomID
	if s.Users[userID].Username != s.Rooms[roomID].Host {
		http.Error(w, "only host can start game", http.StatusUnauthorized)
	}

	game := models.Game{
		Turn:        s.Users[userID].Username,
		TrumpNumber: 2,
	}

	players := make([]*models.Player, len(s.Rooms[roomID].Users))
	for i, u := range s.Rooms[roomID].Users {
		players[i] = &models.Player{
			Username: u,
			Level:    2,
		}
	}

	s.Rooms[roomID].Game = &game
	s.broadcastUpdate(roomID, "game_started")
	w.WriteHeader(http.StatusOK)
}
