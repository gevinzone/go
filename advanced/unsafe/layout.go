package practice

import (
	"fmt"
	"reflect"
)

func PrintFieldOffset(entity any) {
	typ := reflect.TypeOf(entity)
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < typ.NumField(); i++ {
		fd := typ.Field(i)
		fmt.Printf("%s: %d\n", fd.Name, fd.Offset)
	}
}
