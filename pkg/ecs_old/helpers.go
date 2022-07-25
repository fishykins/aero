package ecs2

import (
	"reflect"
	"runtime"
	"strings"
)

func NewQuery(with []string, without []string) Query {
	compiledWith := make([]QueryComponent, len(with))
	for i, component := range with {
		compiledWith[i] = QueryComponent{Type: component}
	}
	return Query{
		With:    compiledWith,
		Without: without,
	}
}

func GetFunctionName(i interface{}) string {
	fullStr := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	return fullStr[strings.LastIndex(fullStr, ".")+1:]
}
