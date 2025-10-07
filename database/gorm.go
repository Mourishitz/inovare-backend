package database

import (
	"log"

	"inovare-backend/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}
