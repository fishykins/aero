package ecscore

import (
	"reflect"
	"runtime"
	"strings"
)

func GetFunctionName(i interface{}) string {
	fullStr := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	return fullStr[strings.LastIndex(fullStr, ".")+1:]
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func NewQuery(tags ...string) Query {
	return Query{
		Components: tags,
	}
}
