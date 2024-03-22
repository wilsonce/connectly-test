package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Access-Token")
		if 1 > len(accessToken) {
			c.JSON(401, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
			return
		}
		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(JwtSecret), nil
		})
		switch {
		case token.Valid:
			tmpClaims := token.Claims.(jwt.MapClaims)
			username := tmpClaims["username"]
			c.Set("username", username)
			c.Next()
		case errors.Is(err, jwt.ErrTokenMalformed):
			c.JSON(http.StatusBadRequest, gin.H{"message": "That's not even a token"})
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			// Invalid signature
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid signature"})
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			c.JSON(http.StatusBadRequest, gin.H{"message": "Timing is everything"})
		default:
			fmt.Println("Couldn't handle this token:", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Couldn't handle this token:%s", err.Error())})
		}
		c.Abort()
	}
}
