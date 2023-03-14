package models

import "gorm.io/gorm"

type Features struct {
	gorm.Model
	Associate   string
	AssociateId int
	IconKey     string
	Value       string
}
