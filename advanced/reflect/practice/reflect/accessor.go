package reflect

import (
	"errors"
	"reflect"
)

type Accessor struct {
	val reflect.Value
	typ reflect.Type
}

func NewReflectAccessor(val any) (*Accessor, error) {
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errors.New("invalid entity")
	}
	return &Accessor{
		typ: reflect.TypeOf(val).Elem(),
		val: reflect.ValueOf(val).Elem(),
	}, nil
}

func (r *Accessor) Field(field string) (int, error) {
	for i := 0; i < r.typ.NumField(); i++ {
		if r.typ.Field(i).Name == field {
			return r.val.Field(i).Interface().(int), nil
		}
	}
	return 0, errors.New("invalid field")
}

func (r *Accessor) FieldV2(field string) (int, error) {
	if _, ok := r.typ.FieldByName(field); !ok {
		return 0, errors.New("invalid field")
	}
	return r.val.FieldByName(field).Interface().(int), nil
}

func (r *Accessor) SetField(field string, val int) error {
	if _, ok := r.typ.FieldByName(field); !ok {
		return errors.New("invalid field")
	}

	fieldVal := r.val.FieldByName(field)
	if !fieldVal.CanSet() {
		return errors.New("unable to set value")
	}
	fieldVal.Set(reflect.ValueOf(val))
	return nil
}
