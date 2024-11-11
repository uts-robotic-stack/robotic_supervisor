package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/yaml.v3"
)

type UserHandler struct {
	users map[string]User
}

// User represents the user information.
type User struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Role     string `json:"role" yaml:"role"`
}

// Claims represents the JWT payload.
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJvYm90aWMiLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3MzA2MTc5NDl9.5RHDAgo7eNv9h3FsW0ypWVQkqKlAiWc8U-V83OszR1Y")

// NewUserHandler initializes UserHandler and loads users from a YAML file.
func NewUserHandler() *UserHandler {
	handler := &UserHandler{users: make(map[string]User)}

	err := handler.loadUsersFromYAML("config/user.yaml")
	if err != nil {
		return nil
	}

	return handler
}

// loadUsersFromYAML reads user data from a YAML file.
func (h *UserHandler) loadUsersFromYAML(filePath string) error {
	yamlData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %v", err)
	}

	var usersData struct {
		Users map[string]User `yaml:"users"`
	}
	if err := yaml.Unmarshal(yamlData, &usersData); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	h.users = usersData.Users
	return nil
}

// SignInHandler handles user sign-in and returns a JWT if successful.
func (h *UserHandler) HandleUserSignIn(c *gin.Context) {
	var credentials User
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, err := h.authenticateUser(credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// If the user role is "user", do not return a token
	if user.Role == "user" {
		c.JSON(http.StatusOK, gin.H{"message": "Authenticated as regular user"})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (h *UserHandler) authenticateUser(username, password string) (User, error) {
	// Iterate over all roles in the users map
	for _, user := range h.users {
		// For other roles, check if both username and password match
		if user.Username == username && user.Password == password {
			return user, nil
		}
	}
	// By default return user
	return User{Username: username, Role: "user"}, nil
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
