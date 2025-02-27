package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/snippet-manager/config"
	"github.com/yourusername/snippet-manager/models"
)

type SnippetInput struct {
	Title    string `json:"title" binding:"required"`
	Language string `json:"language" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

func CreateSnippet(c *gin.Context) {
	var input SnippetInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userIdStr, _ := c.Get("userId")
	userIdUint, err := strconv.ParseUint(userIdStr.(string), 10, 32)
	fmt.Println(userIdUint, "idd babyy")

	if err != nil {
		c.JSON(500, gin.H{"error": "Invalid user ID format"})
		return
	}

	var createdSnippet models.Snippet = models.Snippet{
		Title:    input.Title,
		Content:  input.Content,
		Language: input.Language,
		UserID:   uint(userIdUint),
	}

	result := config.DB.Create(&createdSnippet)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Error creating snippet"})
		return
	}

	c.JSON(200, gin.H{"data": createdSnippet})
}

func GetSnippets(c *gin.Context) {
	userIdStr, e := c.Get("userId")

	if !e {
		c.JSON(500, gin.H{"error": "Invalid user ID format"})
		return
	}

	userIdUint, err := strconv.ParseUint(userIdStr.(string), 10, 32)

	if err != nil {
		c.JSON(500, gin.H{"error": "Invalid user ID format"})
		return
	}

	var allSnippets []models.Snippet

	config.DB.Find(&allSnippets, "user_id = ?", userIdUint)
	c.JSON(200, gin.H{"data": allSnippets})
}

func GetSnippet(c *gin.Context) {
	snippetId := c.Param("id")

	// Parse the user ID from the token
	userId, err := ParseUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	var snippet models.Snippet

	// Query the snippet where both user_id and id match
	result := config.DB.Where("user_id = ? AND id = ?", userId, snippetId).First(&snippet)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Snippet not found"})
		return
	}

	// Return the snippet
	c.JSON(http.StatusOK, gin.H{"data": snippet})
}

func UpdateSnippet(c *gin.Context) {
	snippetId := c.Param("id")

	// Parse the user ID from the token
	userId, err := ParseUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	var snippet models.Snippet

	// Query the snippet where both user_id and id match
	result := config.DB.Where("user_id = ? AND id = ?", userId, snippetId).First(&snippet)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Snippet not found"})
		return
	}

	var input SnippetInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	snippet.Title = input.Title
	snippet.Language = input.Language
	snippet.Content = input.Content

	config.DB.Save(&snippet)

	c.JSON(http.StatusOK, gin.H{"data": snippet})
}

func DeleteSnippet(c *gin.Context) { 
	snippetId := c.Param("id")

	// Parse the user ID from the token
	userId, err := ParseUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	var snippet models.Snippet

	// Query the snippet where both user_id and id match
	result := config.DB.Where("user_id = ? AND id = ?", userId, snippetId).First(&snippet)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Snippet not found"})
		return
	}

	config.DB.Delete(&snippet)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func ParseUserId(c *gin.Context) (uint, error) {
	userIdStr, e := c.Get("userId")

	if !e {
		return 0, fmt.Errorf(`invalid user ID format`)
	}

	userIdUint, err := strconv.ParseUint(userIdStr.(string), 10, 32)
	if err != nil {
		return 0, fmt.Errorf(`invalid user ID format`)
	}

	return uint(userIdUint), nil
}
