package server

import (
	"github.com/eycai/tractor/src/internal/models"
)

func (s *Server) createRoom(userID string, name string, capacity int) *models.Room {
	creator := s.Users[userID].Username
	roomID := s.generateRoomID()

	// create room
	s.Rooms[roomID] = &models.Room{
		ID:       roomID,
		Name:     name,
		Users:    []string{},
		Host:     creator,
		Capacity: capacity,
	}

	s.addToRoom(userID, roomID)
	return s.Rooms[roomID]
}

func (s *Server) addToRoom(userID string, roomID string) {
	s.Users[userID].RoomID = roomID
	s.Rooms[roomID].Users = append(s.Rooms[roomID].Users, s.Users[userID].Username)
}

func (s *Server) removeFromRoom(userID string, roomID string) {
	s.Users[userID].RoomID = ""
	user := s.Users[userID].Username
	i := indexOf(s.Rooms[roomID].Users, user)
	s.Rooms[roomID].Users = remove(s.Rooms[roomID].Users, i)
}
