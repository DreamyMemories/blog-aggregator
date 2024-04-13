package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DreamyMemories/blog-aggregator/httpfunctions"
	"github.com/DreamyMemories/blog-aggregator/internal/database"
	"github.com/DreamyMemories/blog-aggregator/types"
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
	apiConfig := &types.ApiConfig{
		DB: dbQueries,
	}

	// Create Server
	mux := httpfunctions.Mux(apiConfig)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	fmt.Println("Hello World", port)
}
