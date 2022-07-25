package ecs

import "log"

func NewQuery(with []string, without []string) Query {
	compiledWith := make([]QueryComponent, len(with))
	for i, component := range with {
		compiledWith[i] = QueryComponent{Type: component}
	}
	ret := Query{
		With:    compiledWith,
		Without: without,
	}

	log.Println("NewQuery:", ret)
	return ret
}
