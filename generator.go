package main

import (
	"gorm.io/gen"
	"rent-n-go-backend/generator"
	"rent-n-go-backend/utils"
)

func generateQuery() {
	var lists generator.Generator

	generator.Generate(&lists)

	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(utils.GetDb())

	g.ApplyBasic(lists.Models...)

	for i, v := range lists.Queries {
		g.ApplyInterface(i, v.([]interface{})...)
	}

	g.Execute()
}
