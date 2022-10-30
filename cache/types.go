package cache

import (
	"encoding/json"
	"errors"
	"github.com/beego/beego/v2/server/web/context"
	"time"
)

var (
	errFailToConvertValueType = errors.New("AnyValue: 无法转换类型")
	ErrKeyNotFound            = errors.New("cache: 找不到 key")
)

type Cache interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, val any, expiration time.Duration) error
	Delete(ctx context.Context, key string) error

	//OnEvicted(func(key string, val any))
}

type AnyValue struct {
	Val any
	Err error
}

func (a AnyValue) String() (string, error) {
	if a.Err != nil {
		return "", a.Err
	}
	str, ok := a.Val.(string)
	if !ok {
		return "", errFailToConvertValueType
	}
	return str, nil
}

func (a AnyValue) BindJson(val any) error {
	if a.Err != nil {
		return a.Err
	}
	str, ok := a.Val.([]byte)
	if !ok {
		return errFailToConvertValueType
	}
	return json.Unmarshal(str, val)
}
