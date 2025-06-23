package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/ramyasreetejo/ShesAmeya/api"
	"github.com/ramyasreetejo/ShesAmeya/model"
)

func main() {
	_ = godotenv.Load()
	apiKey := os.Getenv("GEMINI_API_KEY")
	port := os.Getenv("REST_PORT")
	if port == "" {
		port = "8080"
	}

	client := model.NewGeminiClient(apiKey)
	http.HandleFunc("/api/chat", api.ChatHandler(client))
	http.Handle("/", http.FileServer(http.Dir("static")))

	log.Println("Server running at http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
