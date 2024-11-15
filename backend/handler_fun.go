package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/IsahiRea/discord-bot/backend/internal/database"
)

func handlerTrivia(w http.ResponseWriter, r *http.Request) {

	type triviaMsg struct {
		Category   string   `json:"category"`
		ID         string   `json:"id"`
		Tags       []string `json:"tags"`
		Difficulty string   `json:"difficulty"`
		IsNiche    bool     `json:"isNiche"`
		Question   struct {
			Text string `json:"text"`
		} `json:"question"`
		CorrectAnswer    string   `json:"correctAnswer"`
		IncorrectAnswers []string `json:"incorrectAnswers"`
		Type             string   `json:"type"`
		Regions          []string `json:"regions"`
	}

	resp, err := http.Get("https://the-trivia-api.com/v2/questions?limit=1&types=text_choice")
	if err != nil {
		respondWithError(w, resp.StatusCode, "couldn't fetch the trivia question")
	}
	defer resp.Body.Close()

	var trivia []triviaMsg
	if err := json.NewDecoder(resp.Body).Decode(&trivia); err != nil {
		respondWithError(w, 500, fmt.Sprintf("error reading trivia json: %v", err))
		return
	}

	sendBack := struct {
		Question         string   `json:"question"`
		CorrectAnswer    string   `json:"correct_answer"`
		IncorrectAnswers []string `json:"incorrect_answers"`
	}{
		trivia[0].Question.Text,
		trivia[0].CorrectAnswer,
		trivia[0].IncorrectAnswers,
	}

	respondWithJSON(w, 200, sendBack)

}

func (cfg *apiConfig) handlerStory(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Message string `json:"message"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, 400, "Couldn't decode parameters")
		return
	}

	// TODO Create a new table for stories
	/*
		sendback = struct{
			Story string `json:"story"`
		}
	*/
}

func (cfg *apiConfig) handlerGenImage(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Message string `json:"message"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, 400, "Couldn't decode parameters")
		return
	}

	resp, err := http.Get("https://pollinations.ai/p/" + params.Message)
	if err != nil {
		respondWithError(w, resp.StatusCode, "Couldn't Generate Image")
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

	// Stream the image data directly to the client
	if _, err := io.Copy(w, resp.Body); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error sending image")
		return
	}

}
