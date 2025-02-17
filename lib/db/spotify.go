package db

import (
	"errors"

	"github.com/m-tsuru/mappier-backend/lib/structs"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

func SaveSpotifyAccessToken(db *gorm.DB, token *oauth2.Token, userId string) (error) {
	accessToken := token.AccessToken
	refreshToken := token.RefreshToken
	expiredAt := token.Expiry

	var chk structs.SpotifyToken
	err := db.Where("id = ?", userId).First(&chk).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		new := structs.SpotifyToken{
			ID: userId,
			AccessToken: accessToken,
			RefreshToken: refreshToken,
			ExpiredAt: expiredAt,
		}
		result := db.Create(&new)
		return result.Error
	}


	chk.AccessToken = accessToken
	chk.RefreshToken = refreshToken
	chk.ExpiredAt = expiredAt
	result := db.Save(&chk)
	return result.Error
}

func SaveSpotifyUser(db *gorm.DB, user structs.SpotifyUserRaw) (error) {
	var chk structs.SpotifyUser
	err := db.Where("id = ?", user.ID).First(&chk).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		new := structs.SpotifyUser{
			ID: user.ID,
			DisplayName: user.DisplayName,
			Email: user.Email,
			ImageUrl: user.Images[0].URL,
		}
		result := db.Create(&new)
		return result.Error
	}
	if err != nil {
		return err
	}
	return nil
}
