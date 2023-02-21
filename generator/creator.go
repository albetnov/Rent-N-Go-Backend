package generator

type Generator struct {
	Models  []interface{}
	Queries map[interface{}]interface{}
}

func (g *Generator) addModel(model interface{}) {
	g.Models = append(g.Models, model)
}

func (g *Generator) applyQuery(query interface{}, models ...interface{}) {
	g.Queries[query] = models
}
