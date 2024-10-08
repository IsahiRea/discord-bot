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

	godotenv.Load(".env")

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
	mux.HandleFunc("GET /v1/api/healthz", handlerReadiness)
	mux.HandleFunc("GET /v1/api/error", handlerErr)

	// Auth Handlers
	mux.HandleFunc("POST /v1/api/login", apiCfg.HandlerLogin)
	mux.HandleFunc("GET /v1/api/refresh", apiCfg.HandlerRefresh)
	mux.HandleFunc("GET /v1/api/revoke", apiCfg.HandlerRevoke)

	// User Handlers
	mux.Handle("GET /v1/api/users/{discordID}", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	mux.HandleFunc("POST /v1/api/users", apiCfg.handlerCreateUser)

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
