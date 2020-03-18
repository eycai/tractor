package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	values := map[string]string{"help": "hi", "oh": "good"}
	jsonValue, _ := json.Marshal(values)
	header := w.Header()
	header.Set("Content-Type", "application/json")
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonValue)
}

func initiateServer() {
	http.HandleFunc("/hello", sayHello)
	buildHandler := http.FileServer(http.Dir("./web"))
	http.Handle("/", buildHandler)

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static")))
	http.Handle("/static/", staticHandler)
}

func main() {
	initiateServer()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
