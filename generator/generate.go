package main

import (
	"rent-n-go-backend/models/UserModels"
)

func generateUserModel(lists *Generator) {
	users := []any{UserModels.User{}, UserModels.Nik{}, UserModels.Sim{}, UserModels.RefreshToken{}, UserModels.UserPhoto{}}

	for _, v := range users {
		lists.addModel(v)
	}
}

func generate(lists *Generator) {
	// UserModels models module
	generateUserModel(lists)
}
