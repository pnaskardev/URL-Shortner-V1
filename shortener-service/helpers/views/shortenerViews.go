package views

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShortenerQueueBody struct {
	ID              uuid.UUID `json:"id"`
	UserID          string    `json:"user_id"`
	ShortenedURLKey string    `json:"shortened_url_key"`
	LongURL         string    `json:"long_url"`
}

type ShortURLDB struct {
	gorm.Model
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID          string    `gorm:"type:varchar(255);not null"`
	ShortenedURLKey string    `gorm:"type:varchar(255);uniqueIndex"`
	LongURL         string    `gorm:"type:text"`
	Deleted         bool      `gorm:"default:false;not null"`
	CreatedAt       int64     `gorm:"autoCreateTime:milli"`
	UpdatedAt       int64     `gorm:"autoUpdateTime:milli"`
}

type ShortenRequest struct {
	Url string `json:"url" validate:"url"`
}
