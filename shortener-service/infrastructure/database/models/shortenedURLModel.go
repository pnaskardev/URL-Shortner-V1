package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProcessedEvents struct {
	EventID   string `gorm:"type:text;not null;primaryKey"`
	CreatedAt int64  `gorm:"autoCreateTime:milli"`
}

type URLs struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ShortURLKey string    `gorm:"type:text;not null;unique;index:idx_short_url_key"`
	LongURL     string    `gorm:"type:text;not null"`
	UserID      string    `gorm:"type:text;not null;index:idx_user_id"`
	IsActive    bool      `gorm:"type:boolean;not null;default:true"`
	ExpiresAt   int64     `gorm:"type:bigint;not null"`
	Deleted     bool      `gorm:"default:false;not null"`
	CreatedAt   int64     `gorm:"autoCreateTime:milli"`
	UpdatedAt   int64     `gorm:"autoUpdateTime:milli"`
}
