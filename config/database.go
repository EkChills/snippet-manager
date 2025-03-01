package config

import (
	"fmt"
	"log"

	"github.com/yourusername/snippet-manager/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	// dsn := fmt.Sprintf(
	// 	"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_PASSWORD"),
	// 	os.Getenv("DB_NAME"),
	// 	os.Getenv("DB_PORT"),
	// )

	database, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	fmt.Println("Database connected")
	DB = database
}

func MigrateDatabase() {
	DB.AutoMigrate(&models.User{}, &models.Snippet{})
	fmt.Println("Database migrated!")
}
