package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Pugpaprika21/go-fiber/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserAuthenticationMiddlewareInterface interface {
	SetTokenJWT(user *dto.UserQueryRow) (string, error)
	JWTProtected() fiber.Handler
	setUserToCookle(c *fiber.Ctx, authHeader string, claims jwt.MapClaims)
}

type userAuthentication struct {
	secretKey []byte
}

func NewUserAuthentication() *userAuthentication {
	return &userAuthentication{
		secretKey: []byte(os.Getenv("SECRET_KEY")),
	}
}

func (u *userAuthentication) SetTokenJWT(user *dto.UserQueryRow) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": user,
		"exp":  time.Now().Add(time.Minute).Unix(),
	})

	tokenJWT, err := token.SignedString(u.secretKey)
	if err != nil {
		return "", err
	}
	return tokenJWT, nil
}

func (u *userAuthentication) JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			c.ClearCookie("userJWTData")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization Header"})
		}

		bearerToken := authHeader[7:]
		token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(u.secretKey), nil
		})

		if err != nil {
			c.ClearCookie("userJWTData")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Expired Token"})
		}

		if !token.Valid {
			c.ClearCookie("userJWTData")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Token"})
		}

		claims, _ := token.Claims.(jwt.MapClaims)

		u.setUserToCookle(c, authHeader, claims)
		return c.Next()
	}
}

func (u *userAuthentication) setUserToCookle(c *fiber.Ctx, authHeader string, claims jwt.MapClaims) {
	user := claims["data"]
	userMap, _ := user.(map[string]interface{})

	id, _ := userMap["ID"].(float64)
	username, _ := userMap["Username"].(string)
	userToken, _ := userMap["Token"].(string)
	tokenJWT := authHeader

	c.Cookie(&fiber.Cookie{
		Name:  "userJWTData",
		Value: fmt.Sprintf("%d|%s|%s|%s", int(id), username, userToken, tokenJWT),
	})
}
