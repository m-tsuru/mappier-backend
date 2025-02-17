package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/m-tsuru/mappier-backend/lib/auth"
	"github.com/m-tsuru/mappier-backend/lib/structs"
	"gorm.io/gorm"
)

func SaveSessionId(db *gorm.DB, userId string, sessionId string, expiredAt time.Time) error{
	var chk structs.Session

	err := db.Where("id = ?", userId).First(&chk).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		new := structs.Session{
			ID: userId,
			SessionKey: sessionId,
			ExpiredAt: expiredAt,
		}
		result := db.Create(&new)
		return result.Error
	}

	chk.SessionKey = sessionId
	chk.ExpiredAt = expiredAt
	result := db.Save(&chk)
	return result.Error
}

func GetUserfromSessionId(db *gorm.DB, sessionId string) (*string, error) {
	var rawResult structs.Session
	err := db.Where("session_key = ?", sessionId).Last(&rawResult).Error
	if err != nil {
		return nil, err
	}
	if rawResult.ExpiredAt.Before(time.Now()) {
		return nil, fmt.Errorf("maybe cookie is expired")
	}
	userId := rawResult.ID
	return &userId, nil
}

func GetAccessTokenfromSessionId(db *gorm.DB, sessionId string) (*string, error) {
	// # セッションからストアしている Spotify のアクセストークンを得る
	userId, err := GetUserfromSessionId(db, sessionId)
	if err != nil {
		return nil, err
	}
	// token
	var token structs.SpotifyToken
	err = db.Where("id = ?", userId).Last(&token).Error
	if err != nil {
		return nil, err
	}
	if token.ExpiredAt.Before(time.Now()) {
		rt := token.RefreshToken
		at, exp, err := auth.RefreshAccessToken(rt)
		if err != nil {
			return nil, fmt.Errorf("unable to get refresh access token: %w", err)
		}
		token.AccessToken = *at
		token.ExpiredAt = time.Now().Add((time.Duration(*exp) - 60) * time.Second)
		err = db.Save(&token).Error
		if err != nil {
			return nil, err
		}
	}
	accessToken := token.AccessToken
	return &accessToken, nil
}
