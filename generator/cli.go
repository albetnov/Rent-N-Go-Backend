package main

import (
	"gorm.io/gen"
	"log"
	"os"
	"rent-n-go-backend/utils"
)

func main() {
	var lists Generator

	os.MkdirAll("./query", 0700)

	generate(&lists)

	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(utils.GetDb())

	generateWithQuery(g.ApplyInterface)
	g.ApplyBasic(lists.Models...)

	g.Execute()

	log.Println("Berhasil di generate!")
}
