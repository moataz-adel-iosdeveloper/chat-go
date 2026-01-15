package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey []byte

func InitJWTService() {
	var secretKey = os.Getenv("SUPER_SECRET_KEY")
	jwtSecretKey = []byte(secretKey)
}

func GenerateToken(userID string) []byte {
	clamis := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, clamis)
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return nil
	}
	return []byte(signedToken)
}

// ValidateToken verifies the JWT signature and returns the user_id claim if valid.
func ValidateToken(tokenStr string) (string, error) {
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Ensure token uses HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecretKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := tok.Claims.(jwt.MapClaims); ok && tok.Valid {
		if uid, ok := claims["user_id"].(string); ok {
			return uid, nil
		}
		return "", errors.New("user_id claim missing")
	}
	return "", errors.New("invalid token")
}
