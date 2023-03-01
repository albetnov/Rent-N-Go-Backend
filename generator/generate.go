package main

import (
	"rent-n-go-backend/models/user"
)

func generateUserModel(lists *Generator) {
	users := []any{user.User{}, user.Nik{}, user.Sim{}, user.RefreshToken{}, user.Photo{}}

	for _, v := range users {
		lists.addModel(v)
	}
}

func generate(lists *Generator) {
	// user models module
	generateUserModel(lists)
}
