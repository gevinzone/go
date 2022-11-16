package rpc

import (
	"context"
	"encoding/json"
	"reflect"
)

func InitClientProxy(service Service, p Proxy) error {
	// 做校验，确保它必须是一个指向结构体的指针

	val := reflect.ValueOf(service).Elem()
	typ := reflect.TypeOf(service).Elem()
	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		fieldType := typ.Field(i)
		fieldVal := val.Field(i)
		if !fieldVal.CanSet() {
			continue
		}
		fn := reflect.MakeFunc(fieldType.Type,
			func(args []reflect.Value) (results []reflect.Value) {
				outType := fieldType.Type.Out(0)
				ctx := args[0].Interface().(context.Context)
				arg := args[1].Interface()
				fillResultWithError := func(err error) {
					results = append(results, reflect.Zero(outType))
					results = append(results, reflect.ValueOf(err))
				}
				bs, err := json.Marshal(arg)

				if err != nil {
					//results = append(results, reflect.Zero(outType))
					//results = append(results, reflect.ValueOf(err))
					fillResultWithError(err)
					return
				}

				req := &Request{
					ServiceName: service.Name(),
					MethodName:  fieldType.Name,
					Data:        bs,
				}
				resp, err := p.Invoke(ctx, req)
				if err != nil {
					//results = append(results, reflect.Zero(outType))
					//results = append(results, reflect.ValueOf(err))
					fillResultWithError(err)
					return
				}
				outVal := reflect.New(outType).Interface()
				err = json.Unmarshal(resp.Data, outVal)
				results = append(results, reflect.ValueOf(outVal).Elem())
				appendErr := func(err error) {
					if err != nil {
						results = append(results, reflect.ValueOf(err))
						return
					}
					results = append(results, reflect.Zero(reflect.TypeOf(new(error)).Elem()))
				}
				appendErr(err)

				return
			})
		fieldVal.Set(fn)
	}

	return nil
}

type Service interface {
	Name() string
}
