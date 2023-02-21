package generator

import "rent-n-go-backend/models"

type Querier struct {
}

func Generate(lists *Generator) {
	lists.addModel(&models.User{})
}
