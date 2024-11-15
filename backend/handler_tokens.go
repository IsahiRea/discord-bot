package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/IsahiRea/discord-bot/backend/internal/auth"
	"github.com/IsahiRea/discord-bot/backend/internal/database"
	"github.com/google/uuid"
)

/*
TODO: Refactor Code
  - creating a model in models.go to simplify sendback
  - Find unnecessary code
*/
func (cfg *apiConfig) checkRefreshToken(w http.ResponseWriter, context context.Context, userID uuid.UUID) bool {

	refreshToken, err := cfg.DB.GetRefreshToken(context, userID)
	if err != nil {
		return false
	}

	if !time.Now().After(refreshToken.ExpiresAt) || !refreshToken.RevokedAt.Valid {
		sendBack := struct {
			RefreshToken string `json:"refresh_token"`
		}{
			refreshToken.Token,
		}

		respondWithJSON(w, 200, sendBack)
	}

	return true
}

func (cfg *apiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		DiscordID string `json:"discord_id"`
		ClientID  string `json:"client_id"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if cfg.ClientID != params.ClientID {
		respondWithError(w, 403, "Unauthorized access to login")
		return
	}

	id, err := parseDiscordID(params.DiscordID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid ID: %v", err))
		return
	}

	user, err := cfg.DB.GetUser(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Couldn't get user: %v", err))
		return
	}

	tokenExists := cfg.checkRefreshToken(w, r.Context(), user.ID)
	if tokenExists {
		return
	}

	//Make the tokens to send back
	newRefreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't create refresh token: %v", err))
		return
	}

	err = cfg.DB.CreateRefeshToken(r.Context(), database.CreateRefeshTokenParams{
		Token:     newRefreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't store refresh token: %v", err))
		return
	}

	sendBack := struct {
		RefreshToken string `json:"refresh_token"`
	}{
		newRefreshToken,
	}

	respondWithJSON(w, 200, sendBack)

}

func (cfg *apiConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 403, fmt.Sprintf("Auth error : %v", err))
		return
	}

	refreshTokenDB, err := cfg.DB.GetRefreshTokenByToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Error finding Refresh Token: %v", err))
		return
	}

	if time.Now().After(refreshTokenDB.ExpiresAt) || refreshTokenDB.RevokedAt.Valid {
		respondWithError(w, 401, fmt.Sprintf("Error finding Refresh Token: %v", err))
		return
	}

	// Create Access Token

	user, err := cfg.DB.GetUserByID(r.Context(), refreshTokenDB.UserID)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Couldn't find user: %v", err))
		return
	}

	newAccessToken, err := auth.MakeJWT(user.DiscordUserID, cfg.TokenSecret)
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

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 403, fmt.Sprintf("Auth error : %v", err))
		return
	}

	if err := cfg.DB.RevokeRefreshToken(r.Context(), refreshToken); err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error revoking refresh token: %v", err))
		return
	}

	respondWithJSON(w, 204, struct{}{})
}
