package auth

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

// TODO Create JWT Authorization package
func GetBearerToken(header http.Header) string {

}

func MakeJWT(tokenString string, tokenSecret string, expiresIn time.Duration) (string, error) {

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

}

func MakeRefreshToken() (string, error) {

}
