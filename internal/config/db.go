package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"todo-list-api/internal/models"
)

func ConnectDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_CONFIG_SERVER")), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	err = db.AutoMigrate(models.Note{}, models.NoteChild{}, models.NoteFile{}, models.NoteChildFiles{})
	if err != nil {
		panic(err.Error())
	}

	log.Printf("Success connecting to database")

	return db
}
