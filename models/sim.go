package models

import "gorm.io/gorm"

type Sim struct {
	gorm.Model
	UserID     uint
	FilePath   string
	IsVerified bool
}
