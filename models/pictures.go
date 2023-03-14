package models

import "gorm.io/gorm"

type Pictures struct {
	gorm.Model
	Associate   string
	AssociateId int
	FileName    string
}
