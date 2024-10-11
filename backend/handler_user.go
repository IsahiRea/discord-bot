package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IsahiRea/discord-bot/backend/internal/database"
)

func (cfg *apiConfig) checkUserExists(context context.Context, discordID int64) bool {

	_, err := cfg.DB.GetUser(context, discordID)

	return err == nil
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		DiscordID string `json:"discord_user_id"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parson JSON: %v", err))
		return
	}

	id, err := parseDiscordID(params.DiscordID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid ID: %v", err))
		return
	}

	userExists := cfg.checkUserExists(r.Context(), id)
	if userExists {
		respondWithJSON(w, 200, struct{}{})
		return
	}

	if err := cfg.DB.CreateUser(r.Context(), id); err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error storing user: %v", err))
		return
	}

	respondWithJSON(w, 201, struct{}{})
}

func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	discordIDSTR := r.PathValue("discordID")

	id, err := parseDiscordID(discordIDSTR)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid ID: %v", err))
		return
	}

	if id != user.DiscordUserID {
		respondWithError(w, 401, fmt.Sprintf("Unauthorized access: %v", err))
		return
	}

	respondWithJSON(w, 200, user)
}
