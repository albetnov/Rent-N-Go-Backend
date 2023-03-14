package models

import "gorm.io/gorm"

type Tour struct {
	gorm.Model
	Name     string
	Desc     string
	Price    int
	Duration string
	Stock    int
	CarId    uint
	DriverId uint
	Car      Cars
	Driver   Driver
}
