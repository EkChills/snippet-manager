package main

import (
	"github.com/yourusername/snippet-manager/config"
	"github.com/yourusername/snippet-manager/routes"
)

func main() {
	config.ConnectDatabase()
	config.MigrateDatabase()


	server := routes.RegisterRoutes()

	server.Run("0.0.0.0:8080")
}