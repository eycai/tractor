package server

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/eycai/tractor/src/internal/api"
	"github.com/eycai/tractor/src/internal/models"
	socketio "github.com/googollee/go-socket.io"
)

type Server struct {
	WSServer          *socketio.Server
	Sockets           map[string]socketio.Conn // map of socket ID to connection
	SocketUsers       map[string]string        // map of socket ID to user ID
	Heartbeats        map[string]*Heartbeat    // map of user ID to last heartbeat
	UserIDLength      int
	RoomIDLength      int
	UserIDs           map[string]string       // map of username to user ID
	Users             map[string]*models.User // map of user ID to user
	Rooms             map[string]*models.Room // map of room id to room
	heartbeatTimeout  time.Duration
	disconnectTimeout time.Duration
	mu                sync.Mutex
}

type Heartbeat struct {
	LastHeartbeat  time.Time
	Disconnected   bool
	PreviousRoomID string
}

func (s *Server) handle(route string, handler http.Handler) {
	handler = api.CORSMiddleware(handler)
	http.Handle(route, handler)
}

func (s *Server) toHandler(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(handler)
}

func (s *Server) handleFunc(route string, handlerFunc http.HandlerFunc) {
	handler := s.toHandler(handlerFunc)
	s.handle(route, handler)
}

func (s *Server) handleRoutes() {
	s.handleFunc("/api/register", s.RegisterUser)
	s.handleFunc("/api/connect", s.ConnectUser)
	s.handleFunc("/api/room_list", s.GetRooms)
	s.handleFunc("/api/join_room", s.JoinRoom)
	s.handleFunc("/api/create_room", s.CreateRoom)
	s.handleFunc("/api/start_game", s.StartGame)
	s.handleFunc("/api/whoami", s.GetUser)
	s.handleFunc("/api/room_info", s.RoomInfo)
	s.handleFunc("/api/heartbeat", s.Heartbeat)
}

func (s *Server) serveClient() {
	buildHandler := http.FileServer(http.Dir("./web"))
	http.Handle("/", buildHandler)
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static")))
	http.Handle("/static/", staticHandler)
}

func (s *Server) handleHeartbeat() {
	for u, h := range s.Heartbeats {
		if h.Disconnected && time.Since(h.LastHeartbeat) > s.disconnectTimeout {
			s.Users[u].Reset()
			delete(s.Heartbeats, u)
		} else if !h.Disconnected && time.Since(h.LastHeartbeat) > s.heartbeatTimeout {
			h.Disconnected = true
			s.Sockets[s.Users[u].SocketID].Close()
			delete(s.Sockets, s.Users[u].SocketID)

			roomID := s.Users[u].RoomID
			if roomID == "" {
				continue
			}

			s.removeFromRoom(u, roomID)
			s.broadcastUpdate(roomID, "player_left")
			h.PreviousRoomID = roomID
		}
	}
}

// func (s *Server) AddToWSRoom(socketID string, roomID string) {
// 	c := s.Sockets[socketID]
// 	c.Join(roomID)
// 	c.SetContext("")
// 	log.Printf("rooms: %v", s.WSServer.Rooms("/"))
// 	log.Printf("members: %v", s.WSServer.RoomLen("/", roomID))
// }

// func (s *Server) LeaveWSRoom(socketID string, roomID string) {
// 	c := s.Sockets[socketID]
// 	c.Leave(roomID)
// }

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
	s.SocketUsers = make(map[string]string)
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

func (s *Server) Start() {
	s.CreateWSServer()
	s.UserIDLength = 8
	s.RoomIDLength = 4
	s.Users = make(map[string]*models.User)
	s.Rooms = make(map[string]*models.Room)
	s.UserIDs = make(map[string]string)
	s.heartbeatTimeout = time.Second
	s.disconnectTimeout = time.Minute
	go s.WSServer.Serve()
	defer s.WSServer.Close()

	s.handle("/socket.io/", s.WSServer)

	s.handleRoutes()
	s.serveClient()

	// start server
	log.Println("Serving at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
