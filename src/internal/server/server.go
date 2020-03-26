package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/eycai/tractor/src/internal/api"
	"github.com/eycai/tractor/src/internal/models"
)

type Server struct {
	WSServer *api.WSServer
	IDLength int
	UserIDs  map[string]string       // map of username to user ID
	Users    map[string]*models.User // map of user ID to user
	Rooms    map[string]*models.Room // map of room id to room
	mu       sync.Mutex
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
	s.handleFunc("/api/leave_room", s.LeaveRoom)
	s.handleFunc("/api/create_room", s.CreateRoom)
	s.handleFunc("/api/start_game", s.StartGame)
	s.handleFunc("/api/whoami", s.GetUser)
	s.handleFunc("/api/room_info", s.RoomInfo)
}

func (s *Server) serveClient() {
	buildHandler := http.FileServer(http.Dir("./web"))
	http.Handle("/", buildHandler)
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static")))
	http.Handle("/static/", staticHandler)
}

func (s *Server) Start() {
	ws := api.NewWSServer()
	s.WSServer = ws
	s.IDLength = 8
	s.Users = make(map[string]*models.User)
	s.Rooms = make(map[string]*models.Room)
	s.UserIDs = make(map[string]string)
	go ws.Server.Serve()
	defer ws.Server.Close()

	s.handle("/socket.io/", ws.Server)

	s.handleRoutes()
	s.serveClient()

	// start server
	log.Println("Serving at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
