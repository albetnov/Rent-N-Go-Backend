package main

import (
	"rent-n-go-backend/models"
	"rent-n-go-backend/models/UserModels"
)

func generateUserModel(lists *Generator) {
	users := []any{UserModels.User{}, UserModels.Nik{}, UserModels.Sim{}, UserModels.RefreshToken{}, UserModels.UserPhoto{}}

	for _, v := range users {
		lists.addModel(v)
	}
}

func generateServicesModel(lists *Generator) {
	services := []any{models.Cars{}, models.Driver{}, models.Tour{}}

	for _, v := range services {
		lists.addModel(v)
	}
}

func generateBasicModel(lists *Generator) {
	basics := []any{models.Features{}, models.Pictures{}}

	for _, v := range basics {
		lists.addModel(v)
	}
}

func generate(lists *Generator) {
	// UserModels models module
	generateUserModel(lists)
	generateServicesModel(lists)
	generateBasicModel(lists)
}
