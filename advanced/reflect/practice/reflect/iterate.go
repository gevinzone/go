package reflect

import (
	"errors"
	"fmt"
	"reflect"
)

// Iterate slice, array and string
func Iterate(input any) ([]any, error) {
	val := reflect.ValueOf(input)
	kind := val.Type().Kind()
	if kind != reflect.Array && kind != reflect.Slice && kind != reflect.String {
		return nil, errors.New("invalid type")
	}
	res := make([]any, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		ele := val.Index(i)
		res = append(res, ele.Interface())
	}
	return res, nil
}

func IterateMapV1(input any) ([]any, []any, error) {
	val := reflect.ValueOf(input)
	if val.Kind() != reflect.Map {
		return nil, nil, errors.New("invalid type")
	}
	keys := make([]any, 0, val.Len())
	values := make([]any, 0, val.Len())
	for _, k := range val.MapKeys() {
		keys = append(keys, k.Interface())
		values = append(values, val.MapIndex(k).Interface())
	}
	return keys, values, nil
}

func IterateMapV2(input any) ([]any, []any, error) {
	val := reflect.ValueOf(input)
	if val.Kind() != reflect.Map {
		return nil, nil, errors.New("invalid type")
	}
	keys := make([]any, 0, val.Len())
	values := make([]any, 0, val.Len())
	iter := val.MapRange()
	for iter.Next() {
		keys = append(keys, iter.Key().Interface())
		values = append(values, iter.Value().Interface())
	}
	iter = val.MapRange()
	fmt.Println(iter.Next())
	return keys, values, nil
}
