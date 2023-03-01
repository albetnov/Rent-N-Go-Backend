package main

import "rent-n-go-backend/models"

func generate(lists *Generator) {
	lists.addModel(models.User{})
	lists.addModel(models.Nik{})
	lists.addModel(models.Sim{})
	lists.addModel(models.RefreshToken{})
	lists.addModel(models.UserPhoto{})
}
