package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/snippet-manager/config"
	"github.com/yourusername/snippet-manager/models"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CheckPassword(hashedPassword string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func GenerateJwt(userId uint) (string, error) {
	fmt.Println(userId, "idd babyy 33")
	userIdStr := strconv.FormatUint(uint64(userId), 10)
	claims := jwt.MapClaims{
		"user_id": userIdStr,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(bearerToken, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing claims"})
			c.Abort()
			return
		}

		// Check expiration
		exp, ok := claims["exp"].(float64)
		if !ok || time.Now().Unix() > int64(exp) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		fmt.Println(claims, "claims")

		// Extract user ID correctly
		userId, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
			c.Abort()
			return
		}

		// Set user ID in the context
		c.Set("userId", userId)
		fmt.Println(userId, "idd babyy set")
		c.Next()
	}
}

func Register(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// if err := config.DB.Where("email = ?", input.Email); err.Error == nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"message": "User already exists",
	// 	})
	// 	return
	// }

	hashedPassword, _ := HashPassword(input.Password)
	user := models.User{Email: input.Email, Password: hashedPassword}
	err := config.DB.Create(&user).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User already exists",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Credentials"})
		return
	}

	if !CheckPassword(user.Password, input.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Credentials"})
		return
	}

	token, _ := GenerateJwt(user.ID)
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}
