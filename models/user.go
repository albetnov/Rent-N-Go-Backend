package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Password string
	Role     string
	Email    string
	Nik      Nik
}
