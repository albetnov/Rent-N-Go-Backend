package models

import (
	"gorm.io/gorm"
	"time"
)

type RefreshToken struct {
	gorm.Model
	UserID    uint
	Token     string
	ExpiredAt time.Time
}
