package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IsahiRea/discord-bot/backend/internal/auth"
	"github.com/IsahiRea/discord-bot/backend/internal/database"
)

func (cfg *apiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		DiscordID int64  `json:"discord_id"`
		ClientID  string `json:"client_id"`
	}

	params := parameters{}
	user, err := cfg.DB.GetUser(r.Context(), params.DiscordID)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Couldn't get user: %v", err))
		return
	}

	refreshToken, err := cfg.DB.GetRefreshToken(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Couldn't get refresh token: %v", err))
		return
	}

	if !time.Now().After(refreshToken.ExpiresAt) || !refreshToken.RevokedAt.Valid {
		//TODO: Refactor this
		sendBack := struct {
			RefreshToken string `json:"refresh_token"`
		}{
			refreshToken.Token,
		}

		respondWithJSON(w, 200, sendBack)
		return
	}

	if cfg.ClientID != params.ClientID {
		respondWithError(w, 403, "Unauthorized access to login")
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
