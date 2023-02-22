package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
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

func Paginate(c *fiber.Ctx) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, err := strconv.Atoi(c.Query("page", "1"))
		RecordLog(err)

		if page == 0 {
			page = 1
		}

		pageSize, err := strconv.Atoi(c.Query("page_size"))
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
