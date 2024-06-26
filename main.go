package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DreamyMemories/blog-aggregator/httpfunctions"
	"github.com/DreamyMemories/blog-aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	dbURL := os.Getenv("CONNECTION_STRING")

	// Load database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database : %v", err)
	}
	dbQueries := database.New(db)
	apiConfig := &httpfunctions.ApiConfig{
		DB: dbQueries,
	}

	// Create Server
	mux := httpfunctions.Mux(apiConfig)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	go httpfunctions.StartScraping(dbQueries, 10, time.Minute)
	log.Println("Starting Server")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
