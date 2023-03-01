package user

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	PhotoPath string
	UserID    uint
}
