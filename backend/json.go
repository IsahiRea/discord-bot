package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with a 5XX error:", msg)
	}

	type errResponse struct {
		Error string `json:"error"`
	}

	errors := errResponse{
		Error: msg,
	}

	respondWithJSON(w, code, errors)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func parseDiscordID(discordID string) (int64, error) {

	id, err := strconv.ParseInt(discordID, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}
