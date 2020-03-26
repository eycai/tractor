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
	defer s.mu.Unlock()

	req := models.JoinRoomRequest{}
	userID, err := s.processRequest(w, r, &req)
	if err != nil {
		return
	}

	// check capacity
	if len(s.Rooms[req.RoomID].Users) >= s.Rooms[req.RoomID].Capacity {
		http.Error(w, "room at capacity", http.StatusBadRequest)
		return
	}

	if s.Rooms[req.RoomID].Game != nil {
		http.Error(w, "game in progress", http.StatusBadRequest)
		return
	}

	s.WSServer.AddToRoom(s.Users[userID].SocketID, req.RoomID)
	s.addToRoom(userID, req.RoomID)
	s.broadcastUpdate(req.RoomID, "update")
	returnSuccess(w)
}

func (s *Server) LeaveRoom(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	req := models.LeaveRoomRequest{}
	userID, err := s.processRequest(w, r, &req)
	if err != nil {
		return
	}

	if s.Rooms[req.RoomID].Game != nil {
		http.Error(w, "game in progress", http.StatusBadRequest)
		return
	}

	s.WSServer.LeaveRoom(s.Users[userID].SocketID, req.RoomID)
	s.removeFromRoom(userID, req.RoomID)
	s.broadcastUpdate(req.RoomID, "update")
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

	req := models.RoomInfoRequest{}
	_, err := s.processRequest(w, r, &req)
	if err != nil {
		return
	}

	room, err := json.Marshal(s.Rooms[req.RoomID])
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
	userID, err := s.processRequest(w, r, &req)
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
	userID, err := s.processRequest(w, r, &req)
	if err != nil {
		return
	}

	roomID := s.createRoom(userID, req.Name, req.Capacity)
	s.WSServer.AddToRoom(s.Users[userID].SocketID, roomID)
	s.emitUpdateToUser(userID, "room_created")

	returnSuccess(w)
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
	s.broadcastUpdate(roomID, "update")
	w.WriteHeader(http.StatusOK)
}

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

func (s *Server) processRequest(w http.ResponseWriter, r *http.Request, req interface{}) (string, error) {
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

func (s *Server) broadcastUpdate(roomID string, eventName string) {
	for _, user := range s.Rooms[roomID].Users {
		userID := s.UserIDs[user]
		s.emitUpdateToUser(userID, eventName)
	}
}

func (s *Server) emitUpdateToUser(userID string, eventName string) {
	update := models.UpdateEvent{
		User: s.Users[userID],
	}

	// update room if it exists
	if s.Users[userID].RoomID != "" {
		if room, ok := s.Rooms[s.Users[userID].RoomID]; ok {
			update.Room = room
		}
	}

	s.emitWSToUser(userID, eventName, update)
}

func (s *Server) emitWSToUser(userID string, eventName string, event interface{}) {
	s.WSServer.Emit(s.Users[userID].SocketID, eventName, event)
}

func returnSuccess(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	resp, err := json.Marshal(&models.EmptyResponse{})
	if err != nil {
		http.Error(w, "error writing response", http.StatusInternalServerError)
	}
	w.Write(resp)
}
