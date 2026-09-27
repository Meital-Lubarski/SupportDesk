package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/conversations", conversationsHandler)
	mux.HandleFunc("/api/conversations/export", conversationsExportHandler)
	mux.HandleFunc("/api/conversations/", conversationByIDHandler)
	mux.HandleFunc("/api/dashboard", dashboardHandler)
	mux.HandleFunc("/api/agents", agentsHandler)

	log.Println("Backend is running at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", enableCORS(mux))) //Start rhe server
}
