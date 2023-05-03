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
	Pictures []Pictures `gorm:"foreignKey:AssociateId"`
}
