package models

import (
	"gorm.io/gorm"
	"rent-n-go-backend/models/UserModels"
	"time"
)

type Orders struct {
	gorm.Model
	CarId         *uint
	DriverId      *uint
	TourId        *uint
	Car           Cars
	Driver        Driver
	Tour          Tour
	TotalAmount   int
	StartPeriod   time.Time
	EndPeriod     time.Time
	UserId        uint
	User          UserModels.User
	PaymentMethod string
	Status        string
}
