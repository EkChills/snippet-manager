package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/snippet-manager/controllers"
)

func RegisterRoutes() *gin.Engine {
	c := gin.Default()

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
