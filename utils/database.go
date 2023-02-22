package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// SatisfiesDbConnection /*
func SatisfiesDbConnection() {
	dbUser := viper.GetString("DB_USER")
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetInt32("DB_PORT")
	dbPass := viper.GetString("DB_PASS")
	dbName := viper.GetString("DB_NAME")

	var credentials string

	if dbPass != "" {
		credentials = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
			dbUser,
			dbPass,
			dbHost,
			dbPort,
			dbName,
		)
	} else {
		credentials = fmt.Sprintf(
			"%s@tcp(%s:%d)/%s?charset=utf8mb4",
			dbUser,
			dbHost,
			dbPort,
			dbName,
		)
	}

	var err error

	db, err = gorm.Open(mysql.Open(credentials), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}

// GetDb
// Return the database instance of Gorm.
func GetDb() *gorm.DB {
	return db
}
