package main

import (
	"fmt"
	"net/http"

	"github.com/IsahiRea/discord-bot/backend/internal/auth"
	"github.com/IsahiRea/discord-bot/backend/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		discordId, err := auth.ValidateJWT(tokenString, cfg.TokenSecret)
		if err != nil {
			respondWithError(w, 401, fmt.Sprintf("Invalid JWT: %v", err))
			return
		}

		user, err := cfg.DB.GetUser(r.Context(), discordId)
		if err != nil {
			respondWithError(w, 404, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
