package practice

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

type FieldAccessor interface {
	Field(field string) (int, error)
	SetField(field string, val int) error
}

type FieldAccessor2 interface {
	Field2(field string) (any, error)
	SetField2(field string, val any) error
}

type FieldMeta struct {
	// offset 后期在我们考虑组合，或者复杂类型字段的时候，它的含义衍生为表达相当于最外层的结构体的偏移量
	offset uintptr
	typ    reflect.Type
}

type UnsafeAccessor struct {
	fields     map[string]FieldMeta
	entityAddr unsafe.Pointer
}

func (u *UnsafeAccessor) Field2(field string) (any, error) {
	ptr, err := u.getFieldPtr(field)
	if err != nil {
		return 0, err
	}
	res := reflect.NewAt(u.fields[field].typ, ptr).Elem().Interface()
	return res, nil
}

func (u *UnsafeAccessor) SetField2(field string, val any) error {
	ptr, err := u.getFieldPtr(field)
	if err != nil {
		return err
	}
	reflect.NewAt(u.fields[field].typ, ptr).Elem().Set(reflect.ValueOf(val))
	return nil
}

func (u *UnsafeAccessor) Field(field string) (int, error) {
	ptr, err := u.getFieldPtr(field)
	if err != nil {
		return 0, err
	}
	res := *(*int)(ptr)
	return res, nil
}

func (u *UnsafeAccessor) SetField(field string, val int) error {
	ptr, err := u.getFieldPtr(field)
	if err != nil {
		return err
	}
	*(*int)(ptr) = val
	return nil
}

func (u *UnsafeAccessor) getFieldPtr(field string) (unsafe.Pointer, error) {
	fdMeta, ok := u.fields[field]
	if !ok {
		return nil, fmt.Errorf("invalid field %s", field)
	}
	ptr := unsafe.Pointer(uintptr(u.entityAddr) + fdMeta.offset)
	if ptr == nil {
		return nil, fmt.Errorf("invalid address of the field: %s", field)
	}
	return ptr, nil
}

func NewUnsafeAccessor(entity any) (*UnsafeAccessor, error) {
	if entity == nil {
		return nil, errors.New("invalid entity")
	}
	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errors.New("invalid entity")
	}
	fields := make(map[string]FieldMeta, typ.Elem().NumField())
	for i := 0; i < typ.Elem().NumField(); i++ {
		fd := typ.Elem().Field(i)
		fields[fd.Name] = FieldMeta{offset: fd.Offset, typ: fd.Type}
	}
	return &UnsafeAccessor{fields: fields, entityAddr: val.UnsafePointer()}, nil
}
