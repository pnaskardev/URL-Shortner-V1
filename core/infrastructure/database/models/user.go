package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username  string    `gorm:"type:varchar(255);not null;unique"`
	Password  string    `gorm:"type:text;not null"`
	CreatedAt int64     `gorm:"autoCreateTime:milli"`
	UpdatedAt int64     `gorm:"autoUpdateTime:milli"`
}
