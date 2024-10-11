package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handlerTrivia(w http.ResponseWriter, r *http.Request) {
	/*
		Every user should have a different daily question
		Fetch from en external api
		serve the json back to the discord bot
	*/

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
		IncorrectAnswers []string `json:"incorrect answer"`
	}{
		trivia[0].Question.Text,
		trivia[0].CorrectAnswer,
		trivia[0].IncorrectAnswers,
	}

	respondWithJSON(w, 200, sendBack)

}