package main

import (
	"log"
	"net/http"
	"websocket/internal/handlers"
)

func main() {
	mux := routes()

	log.Println("Starting web server on port 8080")

	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()

	_ = http.ListenAndServe(":8080", mux)
}
