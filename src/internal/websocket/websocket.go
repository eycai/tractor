package websocket

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func NewServer() *socketio.Server {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Printf("connected: %s", s.ID())
		s.Emit("reply", "connection successful")
		return nil
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		log.Printf("error: %v", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Printf("closed: %s", reason)
	})

	return server
}
