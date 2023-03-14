package models

import "gorm.io/gorm"

type Driver struct {
	gorm.Model
	Name         string
	Desc         string
	Price        int
	Availability bool
}
