package db

import (
	"fmt"

	"github.com/m-tsuru/mappier-backend/lib/structs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


func Setup() error {
	db, err := Open()
	if err != nil {
		return err
	}

	err = Migrate(db)
	if err != nil {
		return err
	}

	return nil
}

func Open() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database connection: %w", err)
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	db.AutoMigrate(structs.SpotifyToken{})
	db.AutoMigrate(structs.SpotifyUser{})
	db.AutoMigrate(structs.Session{})
	return nil
}
