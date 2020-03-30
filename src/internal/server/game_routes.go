package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eycai/tractor/src/internal/models"
)

func (s *Server) getRoom(userID string) (*models.Room, error) {
	if room, ok := s.Rooms[s.Users[userID].RoomID]; ok {
		return room, nil
	}
	return nil, fmt.Errorf("no such room")
}

func (s *Server) usernamesToUsers(usernames []string) (map[string]*models.User, error) {
	users := make(map[string]*models.User)
	for _, u := range usernames {
		if user, ok := s.Users[u]; ok {
			users[u] = user
		} else {
			return users, fmt.Errorf("Invalid username")
		}
	}
	return users, nil
}

func (s *Server) BeginDrawing(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	userID := s.getUserID(w, r)
	if userID == "" {
		return
	}

	room, err := s.getRoom(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if room.Game == nil {
		http.Error(w, "no game started", http.StatusConflict)
		return
	}

	game := room.Game
	if s.Users[userID].Username != room.Host {
		http.Error(w, "only host can start game", http.StatusConflict)
		return
	}

	deck := game.GetDeck()

	// username to user
	usernames := room.Usernames()
	users, err := s.usernamesToUsers(usernames)
	if err != nil {
		http.Error(w, "invalid username given", http.StatusConflict)
		return
	}

	go s.dealCards(users, usernames, deck, room)
	returnSuccess(w)
}

func (s *Server) dealCards(users map[string]*models.User, dealOrder []string, deck models.Deck, room *models.Room) {
	dealUser := 0
	for _, c := range deck.Cards {
		s.mu.Lock()
		users[dealOrder[dealUser]].DealCard(c)
		s.mu.Unlock()
		s.emitUpdateToUser(users[dealOrder[dealUser]].ID, "card_drawn")
		log.Printf("draw card %s, %v", dealOrder[dealUser], c)
		dealUser = (dealUser + 1) % len(dealOrder)
	}
	s.mu.Lock()
	room.Game.GamePhase = models.DrawingComplete
	s.mu.Unlock()
	s.broadcastUpdate(room.ID, "drawing_complete")
}

func (s *Server) FlipCards(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	req := models.FlipCardsRequest{}
	userID, err := s.processPostRequest(w, r, &req)
	if err != nil {
		return
	}
	room, err := s.getRoom(userID)
	if err != nil {
		http.Error(w, "no such room", http.StatusConflict)
		return
	}

	if room.Game == nil {
		http.Error(w, "no game attached", http.StatusConflict)
		return
	}

	if ok := room.Game.FlipCard(req.Card, req.NumCards, s.Users[userID].Username); !ok {
		http.Error(w, "invalid flip", http.StatusConflict)
		return
	}

	s.broadcastUpdate(room.ID, "trump_chosen")
	returnSuccess(w)
}
