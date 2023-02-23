package generator

import "rent-n-go-backend/models"

func Generate(lists *Generator) {
	lists.addModel(models.User{})
}
