package main

import (
	"fmt"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"log"
	"os"
	"rent-n-go-backend/models"
	"rent-n-go-backend/models/UserModels"
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

func processMigration(args []string) {
	if len(args) > 0 {
		supportedArgs := []string{"refresh", "migrate", "seed"}

		if !slices.Contains(supportedArgs, args[0]) {
			log.Fatalf("Unsupported argument: %s", args[0])
			os.Exit(1)
		}

		if args[0] == "refresh" {
			fmt.Println("Refreshing migration...")
			tables, err := utils.GetDb().Migrator().GetTables()

			if err != nil {
				panic(err)
			}

			interfaces := make([]interface{}, len(tables))

			for i, v := range tables {
				interfaces[i] = v
			}

			utils.GetDb().Migrator().DropTable(interfaces...)
		}

		if args[0] == "migrate" || args[0] == "refresh" {
			fmt.Println("Migrating database...")
			// migrate all tables to database.
			migrate(utils.GetDb())
		}

		if args[0] == "seed" || args[0] == "refresh" {
			fmt.Println("Seeding database...")
			// check if UserModels want to seed specific module
			arg := ""

			if len(args) > 1 {
				arg = args[1]
			}

			// seed the table based on arguments.
			seeder(utils.GetDb(), arg)
		}

		fmt.Println("Action executed successfully")
		os.Exit(0)
	}
}

func migrateUserModule(db *gorm.DB) {
	usersModels := []any{&UserModels.User{}, &UserModels.Nik{}, &UserModels.Sim{}, &UserModels.RefreshToken{}, &UserModels.UserPhoto{}}

	for _, v := range usersModels {
		migrateModel(db, v)
	}
}

func migrateBasicModule(db *gorm.DB) {
	basicModels := []any{&models.Features{}, &models.Pictures{}}

	for _, v := range basicModels {
		migrateModel(db, v)
	}
}

func migrateServicesModule(db *gorm.DB) {
	serviceModels := []any{&models.Cars{}, &models.Tour{}, &models.Driver{}}

	for _, v := range serviceModels {
		migrateModel(db, v)
	}
}

/*
*
Will be executed by GORM
in Before Hook
*/
func migrate(db *gorm.DB) {
	migrateUserModule(db)
	migrateBasicModule(db)
	migrateServicesModule(db)
	migrateModel(db, &models.Orders{})
}

// Seed a data to a database
// produce fake data that will only be seeded under development state
// Will be executed in Before Hook.
func seeder(db *gorm.DB, args string) {
	seedByModule(args, "UserModels", func() {
		password, _ := utils.HashPassword("1")

		user := UserModels.User{
			Name:        "Sang Admin",
			Email:       "admin@mail.com",
			Role:        "admin",
			Password:    password,
			PhoneNumber: "0829849434",
		}

		db.Create(&user)
	})
}
