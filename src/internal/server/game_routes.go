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
		http.Error(w, "only host can start drawing", http.StatusConflict)
		return
	}

	deck := game.GetDeck()

	// draw order
	usernames := room.DrawOrder()
	users, err := s.usernamesToUsers(usernames)
	if err != nil {
		http.Error(w, "invalid username given", http.StatusConflict)
		return
	}

	room.Game.GamePhase = models.Drawing
	go s.dealCards(users, usernames, deck, room)
	returnSuccess(w)
}

func (s *Server) dealCards(users map[string]*models.User, dealOrder []string, deck models.Deck, room *models.Room) {
	dealUser := 0
	for i := 0; i < len(deck.Cards)-room.Game.KittySize(); i++ {
		s.mu.Lock()
		users[dealOrder[dealUser]].DealCard(deck.Cards[i])
		s.mu.Unlock()
		s.emitUpdateToUser(users[dealOrder[dealUser]].ID, "card_drawn")
		log.Printf("draw card %s, %v", dealOrder[dealUser], deck.Cards[i])
		dealUser = (dealUser + 1) % len(dealOrder)
	}
	s.mu.Lock()
	room.Game.GamePhase = models.DrawingComplete
	room.Game.SetKitty(deck.Cards[len(deck.Cards)-room.Game.KittySize():])
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

	if room.Game.IsFirstRound() {
		room.Game.SetBanker(s.Users[userID].Username)
	}
	s.broadcastUpdate(room.ID, "trump_chosen")
	returnSuccess(w)
}

func (s *Server) GetKitty(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userID := s.getUserID(w, r)
	if userID == "" {
		return
	}

	room, err := s.getRoom(userID)
	if err != nil {
		return
	}
	username := s.Users[userID].Username
	if username != room.Game.Banker {
		http.Error(w, "user is not banker", http.StatusConflict)
		return
	}

	if room.Game.GamePhase != models.DrawingComplete {
		http.Error(w, "cannot get kitty in current phase", http.StatusConflict)
	}

	room.Game.GamePhase = models.SetKitty

	// update cards in hands based on trump
	users, err := s.usernamesToUsers(room.DrawOrder())
	if err != nil {
		http.Error(w, "invalid username", http.StatusConflict)
		return
	}

	vals := room.Game.GetCardValues()
	for _, u := range users {
		hand := make([]models.Card, len(u.Hand))
		for i, c := range u.Hand {
			hand[i] = c.WithGameValues(vals)
			hand[i] = hand[i].WithTrump(room.Game.TrumpNumber, room.Game.TrumpSuit)
		}
		u.Hand = hand
	}

	kitty := room.Game.GetKitty()

	for i, c := range kitty {
		kitty[i] = c.WithGameValues(vals)
		kitty[i] = kitty[i].WithTrump(room.Game.TrumpNumber, room.Game.TrumpSuit)
	}

	users[room.Game.Banker].Kitty = kitty
	s.broadcastUpdate(room.ID, "cards_finalized")
	returnSuccess(w)
}

func (s *Server) SetKitty(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	req := models.SetKittyRequest{}
	userID, err := s.processPostRequest(w, r, &req)
	if err != nil {
		return
	}

	room, err := s.getRoom(userID)
	if err != nil {
		http.Error(w, "room does not exist", http.StatusConflict)
		return
	}

	user := s.Users[userID]
	username := user.Username
	if username != room.Game.Banker {
		http.Error(w, "user is not banker", http.StatusConflict)
		return
	}

	if len(user.Kitty) != room.Game.KittySize() {
		http.Error(w, "kitty is incorrect length", http.StatusConflict)
		return
	}

	// update kitty
	user.Hand = req.Hand
	user.Kitty = req.Kitty

	// reassign values after losing them
	vals := room.Game.GetCardValues()
	hand := make([]models.Card, len(user.Hand))
	for i, c := range user.Hand {
		hand[i] = c.WithGameValues(vals)
		hand[i] = hand[i].WithTrump(room.Game.TrumpNumber, room.Game.TrumpSuit)
	}
	user.Hand = hand
	s.broadcastUpdate(room.ID, "round_started")
	returnSuccess(w)
}

func (s *Server) PlayCards(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	req := models.PlayCardsRequest{}
	userID, err := s.processPostRequest(w, r, &req)
	if err != nil {
		return
	}

	room, err := s.getRoom(userID)
	if err != nil {
		http.Error(w, "room does not exist", http.StatusConflict)
		return
	}

	if !room.Game.IsValidPlay(req.Cards, s.Users[userID].Hand) {
		http.Error(w, "play invalid", http.StatusConflict)
		return
	}

	users, err := s.usernamesToUsers(room.DrawOrder())
	if err != nil {
		http.Error(w, "invalid users", http.StatusConflict)
	}
	hands := [][]models.Card{}
	for _, u := range users {
		hands = append(hands, u.Hand)
	}
	status, err := room.Game.PlayCards(s.Users[userID].Username, req.Cards, hands)
	if err != nil {
		http.Error(w, "invalid play", http.StatusConflict)
	}

	// play cards from hand
	s.Users[userID].PlayCards(req.Cards)
	if status == models.PlayingTrick {
		// set next turn
		nextUserIndex := (indexOf(room.Users, s.Users[userID].Username) + 1) % len(room.Users)
		room.Game.Turn = room.Users[nextUserIndex].Username
		s.broadcastUpdate(room.ID, "cards_played")
		return
	}

	room.Game.EndTrick()
	if len(s.Users[userID].Hand) == 0 {
		// round ended
		room.Game.EndRound(room.DrawOrder())
		s.broadcastUpdate(room.ID, "round_ended")
	} else {
		room.Game.EndTrick()
		s.broadcastUpdate(room.ID, "trick_ended")
	}
	returnSuccess(w)
}
