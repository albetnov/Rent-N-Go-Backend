package main

import (
	"rent-n-go-backend/models"
	"rent-n-go-backend/utils"
)

/*
*
Will be executed by GORM
in Before Hook
*/
func migrate() {
	utils.GetDb().AutoMigrate(&models.User{})
}
