package main

import (
	"gorm.io/gorm"
	"rent-n-go-backend/models"
	"rent-n-go-backend/utils"
)

func migrateModel(db *gorm.DB, userModel any) {
	err := db.AutoMigrate(userModel)

	if err != nil {
		utils.ShouldPanic(err)
	}
}

/*
*
Will be executed by GORM
in Before Hook
*/
func migrate(db *gorm.DB) {
	migrateModel(db, &models.User{})
}
