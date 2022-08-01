package ecs

import (
	"hash/fnv"
	"reflect"
	"runtime"
	"strings"
)

func GetTypeId(i interface{}) string {
	component, ok := i.(Labled)
	if !ok {
		// Just get the type name if it's not a component.
		fullStr := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
		return fullStr[strings.LastIndex(fullStr, ".")+1:]
	} else {
		return component.Type()
	}
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

func HashTags(tags []string) uint32 {
	h := fnv.New32a()
	s := strings.Join(tags[:], ",")
	h.Write([]byte(s))
	return h.Sum32()
}
