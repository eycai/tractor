package server

import (
	"log"

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

	log.Printf("new room: %s", roomID)
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
	room, ok := s.Rooms[roomID]
	if !ok {
		return
	}
	i := indexOf(room.Users, user)
	room.Users = remove(room.Users, i)
	if user == room.Host && len(room.Users) > 0 {
		room.Host = room.Users[0]
	}

	if len(room.Users) == 0 {
		delete(s.Rooms, roomID)
	}
	log.Printf("%v", room.Users)
}
