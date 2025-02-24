package controllers

import (
	"net/http"
	"os"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/snippet-manager/config"
	"github.com/yourusername/snippet-manager/models"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email	string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CheckPassword(hashedPassword string, password string) bool{
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func GenerateJwt(userId string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func Register(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := HashPassword(input.Password)
	user := models.User{Email: input.Email, Password: hashedPassword}
	config.DB.Create(&user)
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}