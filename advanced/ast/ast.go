package ast

import (
	"fmt"
	"go/ast"
	"reflect"
)

type printVisitor struct {
}

func (p *printVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		//fmt.Println(nil)
		return p
	}

	val := reflect.ValueOf(node)
	typ := reflect.TypeOf(node)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}
	fmt.Printf("val: %+v, type: %s \n", val.Interface(), typ.Name())

	return p
}
