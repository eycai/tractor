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
		Users:    []*models.UserStatus{},
		Host:     creator,
		Capacity: capacity,
	}

	log.Printf("new room: %s", roomID)
	s.addToRoom(userID, roomID)
	return s.Rooms[roomID]
}

func (s *Server) addToRoom(userID string, roomID string) {
	s.Users[userID].RoomID = roomID
	s.Rooms[roomID].Users = append(s.Rooms[roomID].Users, &models.UserStatus{
		Username:  s.Users[userID].Username,
		Connected: true,
	})
}

func (s *Server) setUserConnectionStatus(userID string, connected bool) {
	if room, ok := s.Rooms[s.Users[userID].RoomID]; ok {
		for _, u := range room.Users {
			if u.Username == s.Users[userID].Username {
				u.Connected = true
				break
			}
		}
		if connected {
			s.broadcastUpdate(s.Users[userID].RoomID, "player_joined")
		} else {
			s.broadcastUpdate(s.Users[userID].RoomID, "player_left")
		}

	}
}
func (s *Server) removeFromRoom(userID string, roomID string) {
	u := s.Users[userID]
	u.RoomID = ""
	u.Hand = []models.Card{}
	u.Kitty = []models.Card{}
	user := u.Username
	room, ok := s.Rooms[roomID]
	if !ok {
		return
	}
	i := indexOf(room.Users, user)
	if i == -1 {
		return
	}
	room.Users = remove(room.Users, i)
	if user == room.Host && len(room.Users) > 0 {
		room.Host = room.Users[0].Username
	}

	if len(room.Users) == 0 {
		delete(s.Rooms, roomID)
		log.Printf("removing room %s", roomID)
	}
	log.Printf("%v", room.Users)
}
