package ginsert

import (
	"reflect"
	"strings"
)

func InsertStmt2(entity interface{}) (string, []interface{}, error) {

	// val := reflect.ValueOf(entity)
	// typ := val.Type()
	// 检测 entity 是否符合我们的要求
	// 我们只支持有限的几种输入

	// 使用 strings.Builder 来拼接 字符串
	// bd := strings.Builder{}

	// 构造 INSERT INTO XXX，XXX 是你的表名，这里我们直接用结构体名字

	// 遍历所有的字段，构造出来的是 INSERT INTO XXX(col1, col2, col3)
	// 在这个遍历的过程中，你就可以把参数构造出来
	// 如果你打算支持组合，那么这里你要深入解析每一个组合的结构体
	// 并且层层深入进去

	// 拼接 VALUES，达成 INSERT INTO XXX(col1, col2, col3) VALUES

	// 再一次遍历所有的字段，要拼接成 INSERT INTO XXX(col1, col2, col3) VALUES(?,?,?)
	// 注意，在第一次遍历的时候我们就已经拿到了参数的值，所以这里就是简单拼接 ?,?,?

	typ, valid := isNotEmptyAndValid(entity)
	if !valid {
		return "", nil, errInvalidEntity
	}

	numField := typ.NumField()
	fieldNames := make([]string, 0, numField)
	fieldValues := make([]interface{}, 0, numField)
	fieldSet := make(map[string]bool, numField)
	extractFieldNamesAndValues2(entity, &fieldNames, &fieldValues, &fieldSet)

	bd := strings.Builder{}
	bd.WriteString("INSERT INTO `")
	bd.WriteString(typ.Name())
	bd.WriteString("`(")
	for i, name := range fieldNames {
		if i > 0 {
			bd.WriteRune(',')
		}
		bd.WriteRune('`')
		bd.WriteString(name)
		bd.WriteRune('`')
	}
	bd.WriteString(") VALUES(")
	for i, _ := range fieldValues {
		if i > 0 {
			bd.WriteRune(',')
		}
		bd.WriteRune('?')
	}

	bd.WriteString(");")

	return bd.String(), fieldValues, nil
}

func isNotEmptyAndValid(entity interface{}) (reflect.Type, bool) {
	if entity == nil {
		return nil, false
	}

	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		val = val.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return typ, false
	}
	for i := 0; i < val.NumField(); i++ {
		fieldTyp := typ.Field(i).Type
		fieldVal := val.Field(i)
		if fieldVal != reflect.Zero(fieldTyp) {
			return typ, true
		}
	}
	return typ, false
}

func getEntityTypeAndValue2(entity interface{}) (reflect.Type, reflect.Value, error) {
	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		val = val.Elem()
	}
	return typ, val, nil
}

func extractFieldNamesAndValues2(entity interface{}, fieldNames *[]string, fieldValues *[]interface{}, fieldSet *map[string]bool) {
	typ, val, _ := getEntityTypeAndValue2(entity)
	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		if typ.Field(i).Anonymous && typ.Field(i).Type.Kind() != reflect.Pointer {
			extractFieldNamesAndValues2(val.Field(i).Interface(), fieldNames, fieldValues, fieldSet)
			continue
		}

		fieldName := typ.Field(i).Name
		if found := (*fieldSet)[fieldName]; found {
			return
		}
		(*fieldSet)[fieldName] = true
		*fieldNames = append(*fieldNames, fieldName)
		*fieldValues = append(*fieldValues, val.Field(i).Interface())

	}
}

func buildValuePlaceHolder2(num int) string {
	s := make([]string, num, num)
	for i := 0; i < num; i++ {
		s[i] = "?"
	}
	return strings.Join(s, ",")
}
