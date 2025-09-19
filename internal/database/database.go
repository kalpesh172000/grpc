package database

import (
	"log"

	"github.com/kalpesh172000/grpc/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(databaseURL string) {
	var err error
	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully!")

	// Auto-migrate tables
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	
	log.Println("Database migration completed!")
}
