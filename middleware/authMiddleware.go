package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thisgirlElan/jwt_auth/initializers"
	"github.com/thisgirlElan/jwt_auth/models"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get cookie from request header
		tokenString, err := c.Cookie("Authorization")

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "You are not logged in.",
			})
			return
		}

		// Parse JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Extract user ID
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRETJTOKEN")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check exp
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token expired",
				})
			}

			// Find User with the token sub
			var user models.User
			initializers.Db.First(&user, claims["sub"])

			if user.ID == 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "user not found",
				})
			}

			c.Set("user", user)

			// Call next middleware function
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized request.",
			})
		}

	}
}
