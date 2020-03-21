package server

import (
	"log"

	"github.com/eycai/tractor/src/internal/models"
	socketio "github.com/googollee/go-socket.io"
)

type WSServer struct {
	Server *socketio.Server
}

func connect(s socketio.Conn) error {
	log.Printf("connected: %s", s.ID())

	s.Emit("connect", models.ConnectEvent{})
	return nil
}

func NewWSServer() *WSServer {
	server, err := socketio.NewServer(nil)
	ws := WSServer{
		Server: server,
	}
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", connect)

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Printf("error: %v", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Printf("closed: %s", reason)
	})

	return &ws
}
