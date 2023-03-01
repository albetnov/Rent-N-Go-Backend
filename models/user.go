package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string
	Password     string
	Role         string
	Email        string `gorm:"unique"`
	Nik          Nik
	Sim          Sim
	RefreshToken RefreshToken
	PhoneNumber  string `gorm:"unique"`
	Photo        UserPhoto
}
