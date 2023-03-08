package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"strconv"
)

var db *gorm.DB

const PAGING_SIZE = "5"
const PAGE_DEFAULT = "1"

// SatisfiesDbConnection
// Create a connection to MySQL.
func SatisfiesDbConnection() {
	dbUser := viper.GetString("DB_USER")
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetInt32("DB_PORT")
	dbPass := viper.GetString("DB_PASS")
	dbName := viper.GetString("DB_NAME")

	var credentials string

	if dbPass != "" {
		credentials = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
			dbUser,
			dbPass,
			dbHost,
			dbPort,
			dbName,
		)
	} else {
		credentials = fmt.Sprintf(
			"%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
			dbUser,
			dbHost,
			dbPort,
			dbName,
		)
	}

	var err error

	db, err = gorm.Open(mysql.Open(credentials))

	if err != nil {
		panic(err)
	}
}

// GetDb
// Return the database instance of Gorm.
func GetDb() *gorm.DB {
	return db
}

// Paginate
// paging a given resource with Fiber Compatible Context
func Paginate(c *fiber.Ctx) func(db gen.Dao) gen.Dao {
	return func(db gen.Dao) gen.Dao {
		page, err := strconv.Atoi(c.Query("page", PAGE_DEFAULT))
		RecordLog(err)

		if page == 0 {
			page = 1
		}

		pageSize, err := strconv.Atoi(c.Query("page_size", PAGING_SIZE))
		RecordLog(err)

		switch {
		case pageSize > 50:
			pageSize = 50
		case pageSize <= 0:
			pageSize = 15
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// InRandomOrder
// Yet another Database scope utility to gen a random item from given db instance
func InRandomOrder(db *gorm.DB) *gorm.DB {
	return db.Order("RAND()")
}
