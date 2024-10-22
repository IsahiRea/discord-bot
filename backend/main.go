package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/IsahiRea/discord-bot/backend/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB          *database.Queries
	TokenSecret string
	ClientID    string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the homepage!")
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	tokenSecret := os.Getenv("JWT_SECRET")
	if portString == "" {
		log.Fatal("JWT_SECRET is not found in the environment")
	}

	clientID := os.Getenv("clientID")
	if clientID == "" {
		log.Fatal("clientID is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	queries := database.New(conn)

	apiCfg := apiConfig{
		DB:          queries,
		TokenSecret: tokenSecret,
		ClientID:    clientID,
	}

	mux := http.NewServeMux() // Create a new multiplexer

	// Default handlers
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/error", handlerErr)

	// Auth Handlers
	mux.HandleFunc("POST /api/login", apiCfg.HandlerLogin)
	mux.HandleFunc("GET /api/refresh", apiCfg.HandlerRefresh)
	mux.HandleFunc("GET /api/revoke", apiCfg.HandlerRevoke)

	// User Handlers
	mux.HandleFunc("GET /api/users/{discordID}", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	// Fun Handlers
	mux.HandleFunc("GET /api/trivias", handlerTrivia)
	mux.HandleFunc("POST /api/stories", apiCfg.middlewareAuth(apiCfg.handlerStory))
	mux.HandleFunc("POST /api/images", apiCfg.middlewareAuth(apiCfg.handlerGenImage))

	/*
		Routes to Build:
		SoundEffect->Voice Channel
	*/

	//Subrouter
	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", mux))

	server := &http.Server{
		Addr:    ":" + portString,
		Handler: mux,
	}

	log.Printf("Server starting on port %v", portString)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port:", portString)
}
