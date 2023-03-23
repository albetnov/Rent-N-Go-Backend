package models

import "gorm.io/gorm"

type Tour struct {
	gorm.Model
	Name     string
	Desc     string
	Price    int
	Duration int
	Stock    int
	CarId    uint
	DriverId uint
	Car      Cars
	Driver   Driver
	Features []Features `gorm:"foreignKey:AssociateId"`
	Pictures []Pictures `gorm:"foreignKey:AssociateId"`
}
