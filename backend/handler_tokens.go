package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IsahiRea/discord-bot/backend/internal/auth"
	"github.com/IsahiRea/discord-bot/backend/internal/database"
)

func (cfg *apiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request, user database.User) {
	//TODO Authenticate the user (Return access and refresh tokens)
}

func (cfg *apiConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request, user database.User) {

	refreshToken, err := cfg.DB.GetRefreshToken(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Error finding Refresh Token: %v", err))
		return
	}

	if time.Now().After(refreshToken.ExpiresAt) || refreshToken.RevokedAt.Valid {
		respondWithError(w, 401, fmt.Sprintf("Error finding Refresh Token: %v", err))
		return
	}

	timeDurationJWT, err := time.ParseDuration("1h")
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error finding Refresh Token: %v", err))
		return
	}

	newAccessToken, err := auth.MakeJWT(user.DiscordUserID, cfg.TokenSecret, timeDurationJWT)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating JWT: %v", err))
		return
	}

	sendBack := struct {
		Token string `json:"access_token"`
	}{
		newAccessToken,
	}

	respondWithJSON(w, 200, sendBack)

}

func (cfg *apiConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	//TODO Revoke the Refresh Token
}
