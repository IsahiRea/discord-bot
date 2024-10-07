package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		DiscordID int64 `json:"discord_user_id"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parson JSON: %v", err))
		return
	}

	if err := cfg.DB.CreateUser(r.Context(), params.DiscordID); err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error storing user: %v", err))
		return
	}

	respondWithJSON(w, 201, struct{}{})
}

func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	discordIDSTR := r.PathValue("discordID")

	id, err := strconv.ParseInt(discordIDSTR, 10, 64)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid ID: %v", err))
		return
	}

	user, err := cfg.DB.GetUser(r.Context(), id)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't find user: %v", err))
		return
	}

	respondWithJSON(w, 200, user)
}
