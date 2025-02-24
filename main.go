package main

import "github.com/yourusername/snippet-manager/config"

func main() {
	config.ConnectDatabase()
	config.MigrateDatabase()
}