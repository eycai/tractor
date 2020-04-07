package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/eycai/tractor/src/internal/models"
)

func (s *Server) GameState(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	userID := s.getUserID(w, r)
	if userID == "" {
		return
	}

	state := models.GameStateResponse{}
	state.User = s.Users[userID]
	if room, ok := s.Rooms[s.Users[userID].RoomID]; ok {
		state.Room = room
	}

	resp, err := json.Marshal(state)
	if err != nil {
		http.Error(w, "error marshalling rooms", http.StatusInternalServerError)
	}

	w.Write(resp)
}

func (s *Server) TestUpdateRoom(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	room := &models.Room{}

	err := json.NewDecoder(r.Body).Decode(room)
	if err != nil {
		http.Error(w, "error decoding request", http.StatusBadRequest)
		return
	}

	log.Printf("room id: %s", room.ID)
	if _, ok := s.Rooms[room.ID]; ok {
		log.Printf("before %v", s.Rooms[room.ID])
		s.Rooms[room.ID] = room
		log.Printf("before %v", s.Rooms[room.ID])
	} else {
		http.Error(w, "no such room", http.StatusBadRequest)
	}
	s.broadcastUpdate(room.ID, "force_update")
	returnSuccess(w)
}

func (s *Server) TestUpdateUser(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := &models.User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, "error decoding request", http.StatusBadRequest)
		return
	}

	log.Printf("user id: %s", user.ID)
	if _, ok := s.Users[user.ID]; ok {
		log.Printf("before %v", s.Users[user.ID])
		s.Users[user.ID] = user
		log.Printf("after %v", s.Users[user.ID])
	} else {
		http.Error(w, "no such user", http.StatusBadRequest)
	}
	returnSuccess(w)
}

func (s *Server) Heartbeat(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	userID := s.getUserID(w, r)
	if userID == "" {
		return
	}

	if h, ok := s.Heartbeats[userID]; ok {
		h.LastHeartbeat = time.Now()
		h.Disconnected = false
	}
}

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

	// s.AddToWSRoom(s.Users[userID].SocketID, req.RoomID)
	if !room.HasUser(s.Users[userID].Username) {
		s.addToRoom(userID, req.RoomID)
		s.broadcastUpdate(req.RoomID, "player_joined")
	}

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

	// check capacity
	if len(room.Users) >= room.Capacity {
		http.Error(w, "room at capacity", http.StatusBadRequest)
		return
	}

	if !room.HasUser(s.Users[userID].Username) {
		returnSuccess(w)
		return
	}

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

	roomID := r.URL.Query().Get("roomId")
	userID := s.getUserID(w, r)
	if userID == "" {
		return
	}

	_, ok := s.Rooms[roomID]
	if !ok {
		http.Error(w, "room does not exist", http.StatusBadRequest)
	}

	roomJSON, err := json.Marshal(s.Rooms[roomID])
	if err != nil {
		http.Error(w, "error marshalling rooms", http.StatusInternalServerError)
	}

	w.Write(roomJSON)
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

	// if already has connection, remove it first
	socket := s.Users[userID].SocketID
	if socket != "" {
		go s.Emit(socket, "stale", models.EmptyResponse{})
	}

	// add socket map
	s.Users[userID].SocketID = req.SocketID
	// s.SocketUsers[req.SocketID] = userID
	log.Printf("updated the socket id to be %v", s.Users[userID].SocketID)
	log.Printf("current users: %v", s.Users)

	s.setUserConnectionStatus(userID, true)

	s.Heartbeats[userID] = &Heartbeat{
		LastHeartbeat: time.Now(),
	}

	returnSuccess(w)
}

func (s *Server) CreateRoom(w http.ResponseWriter, r *http.Request) {
	req := models.CreateRoomRequest{}
	userID, err := s.processPostRequest(w, r, &req)
	if err != nil {
		return
	}

	room := s.createRoom(userID, req.Name, req.Capacity)
	// s.AddToWSRoom(s.Users[userID].SocketID, room.ID)
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
		GamePhase:   models.Start,
	}

	players := make(map[string]*models.Player, len(s.Rooms[roomID].Users))
	for _, u := range s.Rooms[roomID].Users {
		players[u.Username] = &models.Player{
			Username: u.Username,
			Level:    2,
		}
	}

	game.Players = players
	s.Rooms[roomID].Game = &game
	playOrder := s.Rooms[roomID].DrawOrder()
	s.Rooms[roomID].Game.SetDrawOrder(playOrder)
	s.broadcastUpdate(roomID, "game_started")
	w.WriteHeader(http.StatusOK)
}
