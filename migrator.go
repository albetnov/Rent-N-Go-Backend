package main

import (
	"gorm.io/gorm"
	"rent-n-go-backend/models"
	"rent-n-go-backend/utils"
)

// migrateModel
// A simple wrapper around GORM Auto Migrate that will automatically complain if error
// based on app state.
func migrateModel(db *gorm.DB, model any) {
	err := db.AutoMigrate(model)

	if err != nil {
		utils.ShouldPanic(err)
	}
}

func seedByModule(args string, module string, callback func()) {
	if args == "" || args == module {
		callback()
	}
}

/*
*
Will be executed by GORM
in Before Hook
*/
func migrate(db *gorm.DB) {
	migrateModel(db, &models.User{})
	migrateModel(db, &models.Nik{})
}

// Seed a data to a database
// produce fake data that will only be seeded under development state
// Will be executed in Before Hook.
func seeder(db *gorm.DB, args string) {
	seedByModule(args, "user", func() {
		password, _ := utils.HashPassword("admin12345")

		user := models.User{
			Name:     "Sang Admin",
			Email:    "admin@mail.com",
			Role:     "admin",
			Password: password,
		}

		db.Create(&user)
	})
}
