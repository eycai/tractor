package server

import (
	"log"

	"github.com/eycai/tractor/src/internal/models"
	socketio "github.com/googollee/go-socket.io"
)

func (s *Server) Emit(wsID string, eventName string, event interface{}) {
	s.Sockets[wsID].Emit(eventName, event)
}

func (s *Server) connectWS(c socketio.Conn) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Printf("connected: %s", c.ID())
	log.Printf("namespace: %s", c.Namespace())
	c.Emit("connect", models.ConnectEvent{
		SocketID: c.ID(),
	})
	s.Sockets[c.ID()] = c
	return nil
}

func (s *Server) disconnectWS(c socketio.Conn, reason string) {
	// s.mu.Lock()
	// defer s.mu.Unlock()
	log.Printf("disconnected %s, reason %s", c.ID(), reason)
	// userID, ok := s.SocketUsers[c.ID()]
	// if !ok {
	// 	return
	// }
	// delete(s.SocketUsers, c.ID())
	// roomID := s.Users[userID].RoomID
	// if s.Users[userID].SocketID != c.ID() {
	// 	// stale tab
	// 	return
	// }
	// s.removeFromRoom(userID, roomID)
	// s.broadcastUpdate(roomID, "player_left")
}

func (s *Server) CreateWSServer() {
	server, err := socketio.NewServer(nil)
	s.Sockets = make(map[string]socketio.Conn)
	// s.SocketUsers = make(map[string]string)
	s.WSServer = server

	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", s.connectWS)

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Printf("error: %v", e)
	})

	server.OnEvent("/", "disconnect", s.disconnectWS)
	server.OnDisconnect("/", s.disconnectWS)
}
