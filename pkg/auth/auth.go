package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func GenerateToken(username string, userID string) (string, error) {
	loadEnv()

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "your-secret-key"
	}

	tokenDuration, err := time.ParseDuration(os.Getenv("JWT_DURATION"))
	if err != nil {
		tokenDuration = time.Hour * 24
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"id":       userID,
		"exp":      time.Now().Add(tokenDuration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("couldn't sign token: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (string, string, error) {
	loadEnv()

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "your-secret-key"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return "", "", fmt.Errorf("couldn't parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", "", fmt.Errorf("invalid username claim")
		}
		userID, ok := claims["id"].(string)
		if !ok {
			return "", "", fmt.Errorf("invalid user ID claim")
		}

		return username, userID, nil
	}

	return "", "", fmt.Errorf("invalid token")
}
