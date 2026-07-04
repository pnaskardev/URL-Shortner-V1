package views

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ShortURLDB is the persistence model for a shortened URL. It stays hand-written
// (not generated from proto) because it carries GORM embedding and column tags.
// The wire contract lives in the contracts module as urlshortenerv1.UrlCreatedEvent.
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
