package models

import "gorm.io/gorm"

type UserPhoto struct {
	gorm.Model
	PhotoPath string
	UserID    uint
}
