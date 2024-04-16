package httpfunctions

import (
	"context"
	"net/http"
	"strings"

	"github.com/DreamyMemories/blog-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiConfig *ApiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		authHeader := r.Header.Get("Authorization")
		const prefix = "ApiKey "
		var apiKey string
		if strings.HasPrefix(authHeader, prefix) {
			apiKey = strings.TrimPrefix(authHeader, prefix)
		} else {
			respondWithError(w, http.StatusUnauthorized, "Incorrect format")
			return
		}

		user, err := apiConfig.DB.GetUserByApiKey(ctx, apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorised API Key detected")
			return
		}
		handler(w, r, user)
	}
}
