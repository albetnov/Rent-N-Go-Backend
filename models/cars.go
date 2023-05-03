package models

import "gorm.io/gorm"

type Cars struct {
	gorm.Model
	Name     string
	Stock    int
	Desc     string
	Price    int
	Pictures []Pictures `gorm:"foreignKey:AssociateId"`
	Seats    int
	Baggage  int
}
