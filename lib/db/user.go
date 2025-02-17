package db

import (
	"github.com/m-tsuru/mappier-backend/lib/structs"
	"gorm.io/gorm"
)

func GetUser(db *gorm.DB, userId string) (*structs.SpotifyUser, error) {
	var user structs.SpotifyUser
	err := db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
