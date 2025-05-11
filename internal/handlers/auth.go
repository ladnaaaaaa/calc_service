package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/ladnaaaaaa/calc_service/internal/database"
	"github.com/ladnaaaaaa/calc_service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type CredentialsRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req CredentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if result := database.DB.Where("login = ?", req.Login).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Create new user
	user := models.User{
		Login:    req.Login,
		Password: string(hashedPassword),
	}

	if result := database.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var req CredentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Find user
	var user models.User
	if result := database.DB.Where("login = ?", req.Login).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"login":   user.Login,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
