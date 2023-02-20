package utils

import (
	"api-enjor/pkg/models"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type PayloadsClaims struct {
	jwt.StandardClaims
	Sub *models.ModuleProfile `json:"sub"`
}

func GenerateTokenJWT(payload *models.ModuleProfile) (string, error) {
	now := time.Now().UTC()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, PayloadsClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(30 * time.Minute).Unix(),
			IssuedAt:  now.Unix(),
		},
		Sub: payload,
	})

	// claims := token.Claims.(jwt.MapClaims)
	// claims["exp"] = now.Add(60 * time.Second).Unix()
	// claims["iat"] = now.Unix()
	// claims["users"] = payload

	decode, _ := Decode(os.Getenv("JWT_SECRET"))

	mySigningKey := []byte((decode))
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func ValidateToken(token string) (interface{}, error) {
	decode, _ := Decode(os.Getenv("JWT_SECRET"))
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return []byte(decode), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalidate token: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("invalid token claim")
	}

	return claims["sub"], nil
}

func ModuleUser(ctx *fiber.Ctx) *models.ModuleProfile {
	user := ctx.Locals("user").(*models.ModuleProfile)
	return user
}
