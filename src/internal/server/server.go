package server

import (
	"log"
	"net/http"
	"sync"
)

type Server struct {
	WSServer     *WSServer
	UserIDLength int
	UserIDs      map[string]string // map of username to user ID
	SocketIDs    map[string]string // map of user ID to socket ID
	mu           sync.Mutex
}

func (s *Server) handle(route string, handler http.Handler) {
	handler = CORSMiddleware(handler)
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
	s.handleFunc("/register", s.RegisterUser)
	s.handleFunc("/connect", s.ConnectUser)
}

func (s *Server) serveClient() {
	buildHandler := http.FileServer(http.Dir("./web"))
	http.Handle("/", buildHandler)
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static")))
	http.Handle("/static/", staticHandler)
}

func (s *Server) Start() {
	ws := NewWSServer()
	s.WSServer = ws
	s.UserIDLength = 8
	s.SocketIDs = make(map[string]string)
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
