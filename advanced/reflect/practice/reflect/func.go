package reflect

import (
	"errors"
	"reflect"
)

func IterateStructFunc(val any) (map[string]*FuncInfo, error) {
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Struct && typ.Kind() != reflect.Ptr {
		return nil, errors.New("invalid type")
	}
	result := make(map[string]*FuncInfo, typ.NumMethod())
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		numIn := method.Type.NumIn()
		ps := make([]reflect.Value, 0, numIn)
		in := make([]reflect.Type, 0, numIn)
		ps = append(ps, reflect.ValueOf(val))
		for i := 0; i < numIn; i++ {
			paramType := method.Type.In(i)
			in = append(in, paramType)
			if i > 0 {
				ps = append(ps, reflect.Zero(paramType))
			}
		}
		ret := method.Func.Call(ps)
		numOut := method.Type.NumOut()
		out := make([]reflect.Type, 0, numOut)
		res := make([]any, 0, numOut)
		for i := 0; i < numOut; i++ {
			out = append(out, method.Type.Out(i))
			res = append(res, ret[i].Interface())
		}

		name := method.Name
		result[name] = &FuncInfo{
			Name:   name,
			In:     in,
			Out:    out,
			Result: res,
		}
	}
	return result, nil
}

type FuncInfo struct {
	Name   string
	In     []reflect.Type
	Out    []reflect.Type
	Result []any
}
