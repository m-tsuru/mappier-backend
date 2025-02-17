package structs

import (
	"time"

	"gorm.io/gorm"
)

type SpotifyToken struct {
	gorm.Model
	ID string `gorm:"primaryKey"`
	AccessToken string
	RefreshToken string
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type SpotifyUser struct {
	gorm.Model
	ID string `gorm:"primaryKey"`
	DisplayName string
	Email string `json:"email"`
	ImageUrl string `json:"image"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type SpotifySongCache struct {
	ID string `gorm:"primaryKey"`
	Image string
	Name string
	ArtistsPureString string // ;(セミコロン区切り)
	Album string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Session struct {
	gorm.Model
	ID string `gorm:"primaryKey"`
	SessionKey string
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
