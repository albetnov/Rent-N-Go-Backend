package main

type Generator struct {
	Models []interface{}
}

func (g *Generator) addModel(model interface{}) {
	g.Models = append(g.Models, model)
}
