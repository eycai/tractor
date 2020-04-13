package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

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
		userID := s.UserIDs[u]
		if user, ok := s.Users[userID]; ok {
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
	log.Printf("users map: %v", users)
	for i := 0; i < len(deck.Cards)-room.Game.KittySize(); i++ {
		s.mu.Lock()
		users[dealOrder[dealUser]].DealCard(deck.Cards[i], room.Game)
		s.mu.Unlock()
		s.emitUpdateToUser(users[dealOrder[dealUser]].ID, "card_drawn")
		log.Printf("draw card %s, %v", dealOrder[dealUser], deck.Cards[i])
		dealUser = (dealUser + 1) % len(dealOrder)
		time.Sleep(time.Millisecond * 500)
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
	for _, p := range room.Game.Players {
		p.ResetCards()
	}
	cardsPlayed := make([]models.Card, req.NumCards)
	for i := 0; i < req.NumCards; i++ {
		cardsPlayed[i] = req.Card
	}
	room.Game.Players[s.Users[userID].Username].PlayCards([][]models.Card{cardsPlayed})

	log.Printf("banker: %s", room.Game.Banker)
	log.Printf("trump: %d %s", room.Game.TrumpNumber, room.Game.TrumpSuit)
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

	kitty := room.Game.GetKitty()
	log.Printf("kitty: %v", kitty)

	users[room.Game.Banker].Kitty = kitty
	users[room.Game.Banker].Hand = append(users[room.Game.Banker].Hand, kitty...)

	for _, u := range users {
		u.Hand = room.Game.GetUpdatedCards(u.Hand)
	}
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

	user.UpdateWithKitty(req.Kitty)

	// reassign values
	user.Hand = room.Game.GetUpdatedCards(user.Hand)

	room.Game.GamePhase = models.Playing
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

	log.Printf("Play: %v", req.Cards)

	room, err := s.getRoom(userID)
	if err != nil {
		log.Printf("room does not exist")
		http.Error(w, "room does not exist", http.StatusConflict)
		return
	}

	if !room.Game.IsValidPlayForGame(req.Cards, s.Users[userID].Hand) {
		log.Printf("play invalid for game")
		http.Error(w, "play invalid", http.StatusConflict)
		return
	}

	users, err := s.usernamesToUsers(room.DrawOrder())
	if err != nil {
		log.Printf("invalid users")
		http.Error(w, "invalid users", http.StatusConflict)
	}

	hands := [][]models.Card{}
	for _, u := range users {
		if u.ID != userID {
			hands = append(hands, u.Hand)
		}
	}

	status, played, err := room.Game.PlayCards(s.Users[userID].Username, req.Cards, hands)
	if err != nil {
		log.Printf("invalid play from playcards")
		http.Error(w, "invalid play", http.StatusConflict)
	}

	// play cards from hand
	s.Users[userID].PlayCards(played)
	if status == models.PlayingTrick {
		// set next turn
		nextUserIndex := (indexOf(room.Users, s.Users[userID].Username) + 1) % len(room.Users)
		room.Game.Turn = room.Users[nextUserIndex].Username
		s.broadcastUpdate(room.ID, "cards_played")
		return
	}

	if len(s.Users[userID].Hand) == 0 {
		// round ended
		peasantPoints, kittyPoints, kitty := room.Game.EndState()
		event := models.EndRoundEvent{
			Kitty:       kitty,
			KittyPoints: kittyPoints,
			TotalPoints: peasantPoints + kittyPoints,
		}
		s.broadcastEvent(room.ID, "round_ended", event)
	} else {
		s.broadcastUpdate(room.ID, "trick_ended")
	}
	log.Printf("success")
	returnSuccess(w)
}

func (s *Server) AdvanceRound(w http.ResponseWriter, r *http.Request) {
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

	room.Game.EndRound()
	s.broadcastUpdate(room.ID, "round_ended")
	returnSuccess(w)
}
