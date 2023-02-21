package main

import (
	"gorm.io/gorm"
	"rent-n-go-backend/models"
)

/*
*
Will be executed by GORM
in Before Hook
*/
func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
}
