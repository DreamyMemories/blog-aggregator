package httpfunctions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/DreamyMemories/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ApiConfig struct {
	DB *database.Queries
}

func Mux(apiConfig *ApiConfig) *mux.Router {
	mux := mux.NewRouter()
	mux.Handle("/v1/readiness", corsMiddleware(http.HandlerFunc(handlerReadiness)))
	mux.Handle("/v1/err", corsMiddleware(http.HandlerFunc(handlerError)))
	mux.Handle("/v1/users", corsMiddleware(apiConfig.middlewareAuth(apiConfig.handlerUser)))
	mux.Handle("/v1/feeds", corsMiddleware(apiConfig.middlewareAuth(apiConfig.handlerFeed)))
	mux.Handle("/v1/allfeeds", corsMiddleware(apiConfig.handlerGetAllFeed()))
	mux.Handle("/v1/feed_follows", corsMiddleware(apiConfig.middlewareAuth(apiConfig.handlerCreateFeedFollow)))
	mux.Handle("/v1/feed_follows/{feedFollowID}", corsMiddleware(apiConfig.middlewareAuth(apiConfig.handlerDeleteFeedFollow))).Methods("DELETE")
	mux.Handle("/v1/posts/{limit}", corsMiddleware(apiConfig.middlewareAuth(apiConfig.handlerGetPostsByUser))).Methods("GET")
	return mux
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handlerError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (apiConfig *ApiConfig) handlerUser(w http.ResponseWriter, r *http.Request, user database.User) {
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

		respondWithJson(w, 200, user)
	} else if r.Method == http.MethodGet {
		respondWithJson(w, 200, user)
	}
}

func (apiConfig *ApiConfig) handlerFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	ctx := context.Background()
	if r.Method == http.MethodPost {
		var body struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil || body.Name == "" || body.Url == "" {
			respondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		newFeed := database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      body.Name,
			Url:       body.Url,
			UserID:    user.ID,
		}

		newFeedFollow := database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			FeedID:    newFeed.ID,
			UserID:    user.ID,
		}

		feed, err := apiConfig.DB.CreateFeed(ctx, newFeed)
		feedFollow, errFeedFollow := apiConfig.DB.CreateFeedFollow(ctx, newFeedFollow)
		if err != nil || errFeedFollow != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create feed"+err.Error())
			return
		}
		fmt.Printf("New feed created: %v\n", feed)
		response := struct {
			Feed        interface{} `json:"feed"`
			Feed_follow interface{} `json:"feed_follow"`
		}{
			Feed:        feed,
			Feed_follow: feedFollow,
		}
		respondWithJson(w, 200, response)
	}
}

func (apiConfig *ApiConfig) handlerGetAllFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			ctx := context.Background()
			feeds, err := apiConfig.DB.GetAllFeed(ctx)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Failed to fetch feeds")
				return
			}
			respondWithJson(w, 200, feeds)
		}
	}
}

func (apiConfig *ApiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	ctx := context.Background()
	if r.Method == "POST" {
		var body struct {
			FeedID uuid.UUID `json:"feed_id"`
		}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		feed, err := apiConfig.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    body.FeedID,
			UserID:    user.ID,
		})

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Somethign went wrong creating feed follow")
			return
		}

		respondWithJson(w, http.StatusOK, feed)
	} else if r.Method == "GET" {
		feeds, err := apiConfig.DB.GetFeedFollowsByUser(ctx, user.ID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Can't find feeds with given user id")
		}

		respondWithJson(w, 200, feeds)

	}
}

func (apiConfig *ApiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	parameter := mux.Vars(r)
	ctx := context.Background()
	feedFollowID := parameter["feedFollowID"]

	id, err := uuid.Parse(feedFollowID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing parameter")
	}

	fmt.Println("deleting feed follow ID:", id)

	feeds, err := apiConfig.DB.DeleteFeedFollow(ctx, id)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to delete id "+err.Error())
	}
	respondWithJson(w, 200, feeds)
}

func (apiConfig *ApiConfig) handlerGetPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	parameter := mux.Vars(r)
	limit := parameter["limit"]
	if limit == "" {
		limit = "10"
	}
	lim, err := strconv.Atoi(limit)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not parse limit")
	}

	posts, err := apiConfig.DB.GetPostsByUser(context.Background(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(lim),
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not create post"+err.Error())
	}
	respondWithJson(w, http.StatusOK, posts)
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
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJson(w, status, map[string]string{"error": message})
}
