package reflect

import (
	"errors"
	"fmt"
	"reflect"
)

func IterateFields(val any) {
	res, err := iterateFields(val)

	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range res {
		fmt.Println(k, v)
	}
}

func iterateFields(val any) (map[string]any, error) {
	if val == nil {
		return nil, errors.New("val can not be nil")
	}

	typ := reflect.TypeOf(val)
	refVal := reflect.ValueOf(val)

	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		refVal = refVal.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("非法类型")
	}
	numFields := typ.NumField()
	res := make(map[string]any, numFields)
	for i := 0; i < numFields; i++ {
		fd := typ.Field(i)
		if typ.Field(i).IsExported() {
			res[fd.Name] = refVal.Field(i).Interface()
		} else {
			res[fd.Name] = reflect.Zero(fd.Type).Interface()
		}

	}
	return res, nil
}

func SetField(entity any, field string, newVal any) error {
	val := reflect.ValueOf(entity)
	typ := val.Type()

	// 只能是一级指针，类似 *User
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return errors.New("非法类型")
	}

	typ = typ.Elem()
	val = val.Elem()

	_, found := typ.FieldByName(field)
	if !found {
		return errors.New("字段不存在")
	}
	fieldVal := val.FieldByName(field)
	if !fieldVal.CanSet() {
		return errors.New("不可修改字段")
	}
	fieldVal.Set(reflect.ValueOf(newVal))
	return nil
}
