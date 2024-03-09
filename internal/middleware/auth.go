package middleware

import (
	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware that checks for a valid authentication token.
func AuthMiddleware(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		validToken := token // Replace with your actual valid token

		if token != validToken {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
