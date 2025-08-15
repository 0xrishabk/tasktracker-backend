package middleware

import (
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func parseJWTFromCookie(c *gin.Context) (jwt.MapClaims, error) {
	cookie, err := c.Cookie("access_token")
	if err != nil || cookie == "" {
		return nil, errors.New("missing auth token")
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return nil, errors.New("missing secret key")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := parseJWTFromCookie(c)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if userID, ok := claims["id"].(string); ok {
			c.Set("userID", userID)
			c.Next()
			return
		}

		c.JSON(401, gin.H{"error": "Invalid ID or userID in token"})
		c.Abort()
	}
}

func JWTAuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := parseJWTFromCookie(c)
		if err == nil {
			if userID, ok := claims["id"].(string); ok {
				c.Set("userID", userID)
			}
		}
		c.Next()
	}
}
