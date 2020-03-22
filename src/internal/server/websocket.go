package server

import (
	"log"

	"github.com/eycai/tractor/src/internal/models"
	socketio "github.com/googollee/go-socket.io"
)

type WSServer struct {
	Server    *socketio.Server
	Sockets   map[string]socketio.Conn // map of socket ID to connection
	Broadcast *socketio.Broadcast
}

func (ws *WSServer) connect(s socketio.Conn) error {
	log.Printf("connected: %s", s.ID())
	log.Printf("namespace: %s", s.Namespace())
	s.Emit("connect", models.ConnectEvent{})
	ws.Sockets[s.ID()] = s
	return nil
}

func (ws *WSServer) addToRoom(socketID string, roomID string) {
	s := ws.Sockets[socketID]
	s.Join(roomID)
	log.Printf("rooms: %v", ws.Server.Rooms("/"))
	log.Printf("members: %v", ws.Server.RoomLen("/", roomID))
}

func (ws *WSServer) leaveRoom(socketID string, roomID string) {
	s := ws.Sockets[socketID]
	s.Leave(roomID)
}

func (ws *WSServer) emit(wsID string, eventName string, event interface{}) {
	ws.Sockets[wsID].Emit(eventName, event)
}

func NewWSServer() *WSServer {
	server, err := socketio.NewServer(nil)
	ws := WSServer{
		Server:  server,
		Sockets: make(map[string]socketio.Conn),
	}
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", ws.connect)

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Printf("error: %v", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Printf("closed: %s", reason)
	})

	return &ws
}
