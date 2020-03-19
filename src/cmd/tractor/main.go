package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/eycai/tractor/src/internal/websocket"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	values := map[string]string{"help": "hi"}
	jsonValue, _ := json.Marshal(values)
	header := w.Header()
	header.Set("Content-Type", "application/json")
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonValue)
}

func main() {
	// r := mux.NewRouter()
	// s := r.PathPrefix("/api").Subrouter()
	// s.Use(websocket.CORSMiddleware)
	// test
	http.HandleFunc("/hello", sayHello)

	// socket.io
	server := websocket.NewServer()
	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", websocket.CORSMiddleware(server))

	// frontend
	buildHandler := http.FileServer(http.Dir("./web"))
	http.Handle("/", buildHandler)
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static")))
	http.Handle("/static/", staticHandler)

	// start server
	log.Println("Serving at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
