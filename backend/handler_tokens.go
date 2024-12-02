package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/IsahiRea/discord-bot/backend/internal/auth"
	"github.com/IsahiRea/discord-bot/backend/internal/database"
)

func (cfg *apiConfig) checkRefreshToken(w http.ResponseWriter, refreshToken database.RefreshToken) bool {

	if !time.Now().After(refreshToken.ExpiresAt) || !refreshToken.RevokedAt.Valid {

		sendBack := RefreshToken{
			Token: refreshToken.Token,
		}
		respondWithJSON(w, 200, sendBack)
		return true
	}

	return false
}

func (cfg *apiConfig) checkUser(w http.ResponseWriter, r *http.Request, params TokenParams) (user database.User, ok bool) {
	if cfg.ClientID != params.ClientID {
		respondWithError(w, 403, "Unauthorized access to login")
		return database.User{}, false
	}

	id, err := parseDiscordID(params.DiscordID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid ID: %v", err))
		return database.User{}, false
	}

	user, err = cfg.DB.GetUser(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Couldn't get user: %v", err))
		return database.User{}, false
	}

	return user, true
}

func (cfg *apiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {

	params := TokenParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, ok := cfg.checkUser(w, r, params)
	if !ok {
		return
	}

	refreshToken, err := cfg.DB.GetRefreshToken(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error finding Refresh Token: %v", err))
		return
	}

	tokenExists := cfg.checkRefreshToken(w, refreshToken)
	if tokenExists {
		return
	}

	//Make the tokens to send back
	newRefreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't create refresh token: %v", err))
		return
	}

	//Store Token into Database
	err = cfg.DB.CreateRefeshToken(r.Context(), database.CreateRefeshTokenParams{
		Token:     newRefreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't store refresh token: %v", err))
		return
	}

	sendBack := RefreshToken{
		Token: newRefreshToken,
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

	tokenExists := cfg.checkRefreshToken(w, refreshTokenDB)
	if tokenExists {
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

	sendBack := AccessToken{
		Token: newAccessToken,
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
