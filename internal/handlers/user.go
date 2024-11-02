package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct{}

// User represents the user information.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// In-memory user storage for demonstration purposes.
var users = map[string]User{
	"robotic": {Username: "robotic", Password: "admin", Role: "admin"},
}

// Claims represents the JWT payload.
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJvYm90aWMiLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3MzA2MTc5NDl9.5RHDAgo7eNv9h3FsW0ypWVQkqKlAiWc8U-V83OszR1Y")

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// SignInHandler handles user sign-in and returns a JWT if successful.
func (h *UserHandler) HandleUserSignIn(c *gin.Context) {
	var credentials User
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	fmt.Println(credentials)
	// Authenticate user credentials.
	user, err := h.authenticateUser(credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Create a new JWT token.
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	// Send the token to the client.
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// authenticateUser verifies the user's credentials.
func (h *UserHandler) authenticateUser(username, password string) (User, error) {
	user, exists := users[username]
	if !exists || user.Password != password {
		return User{}, errors.New("invalid credentials")
	}
	return user, nil
}

// RoleHandler returns the role of the authenticated user.
func (h *UserHandler) HandleRole(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": claims.Role})
}
