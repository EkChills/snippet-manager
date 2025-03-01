package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yourusername/snippet-manager/controllers"
)

func RegisterRoutes() *gin.Engine {
	c := gin.Default()

	c.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins (change for security)
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	apiGroup := c.Group("/api/v1")

	snippetGroup := apiGroup.Group("/snippets")
	snippetGroup.Use(controllers.AuthMiddleware())
	

	apiGroup.POST("/register", controllers.Register)
	apiGroup.POST("/login", controllers.Login)
	snippetGroup.POST("/", controllers.CreateSnippet)
	snippetGroup.GET("/", controllers.GetSnippets)
	snippetGroup.GET("/:id", controllers.GetSnippet)
	snippetGroup.PATCH("/:id", controllers.UpdateSnippet)
	snippetGroup.DELETE("/:id", controllers.DeleteSnippet)

	return c

}
