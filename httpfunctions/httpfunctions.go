package httpfunctions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DreamyMemories/blog-aggregator/internal/database"
	"github.com/DreamyMemories/blog-aggregator/types"
	"github.com/google/uuid"
)

func Mux(apiConfig *types.ApiConfig) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/v1/readiness", corsMiddleware(http.HandlerFunc(handlerReadiness)))
	mux.Handle("/v1/err", corsMiddleware(http.HandlerFunc(handlerError)))
	mux.Handle("/v1/users", corsMiddleware(handlerUser(apiConfig)))
	return mux
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handlerError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

func handlerUser(apiConfig *types.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			ctx := context.Background()
			var name struct {
				Name string `json:"name"`
			}
			err := json.NewDecoder(r.Body).Decode(&name)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid request body")
				return
			}

			newUser := database.CreateUserParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      name.Name,
			}

			// Create user in dbqueries
			user, err := apiConfig.DB.CreateUser(ctx, newUser)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Failed to create user "+err.Error())
			}
			fmt.Println(user)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
		} else if r.Method == http.MethodGet {
			ctx := context.Background()
			authHeader := r.Header.Get("Authorization")
			const prefix = "ApiKey "
			var apiKey string
			if strings.HasPrefix(authHeader, prefix) {
				apiKey = strings.TrimPrefix(authHeader, prefix)
			} else {
				respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			}

			user, err := apiConfig.DB.GetUserByApiKey(ctx, apiKey)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Unauthorised API Key detected")
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
		}
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Checks if request is CORS preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		//Call next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func respondWithJson(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
	return
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJson(w, status, map[string]string{"error": message})
}
