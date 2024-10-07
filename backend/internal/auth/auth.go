package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO Create JWT Authorization package
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no authentication info found")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 {
		return "", errors.New("malformed auth header")
	}
	if tokenParts[0] != "Bearer" {
		return "", errors.New("malformed first part of auth header")
	}

	return tokenParts[1], nil
}

func MakeJWT(userID int64, tokenSecret string, expiresIn time.Duration) (string, error) {

	mySigningKey := []byte(tokenSecret)

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   strconv.FormatInt(userID, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "nil", err
	}

	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("error parsing claims: %s", err)
	}

	discordIDSTR, err := token.Claims.GetSubject()
	if err != nil {
		return 0, fmt.Errorf("error obtaining discordId from claims: %s", err)
	}

	discordId, err := strconv.ParseInt(discordIDSTR, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ID: %s", err)
	}

	return discordId, nil
}

func MakeRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("error generating bytes: %v", err)
	}

	return hex.EncodeToString(randomBytes), nil
}
