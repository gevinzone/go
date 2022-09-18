package sqlp

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonColumn_Value(t *testing.T) {
	jc := JsonColumn[User]{Valid: true, Val: User{Name: "Tom"}}
	value, err := jc.Value()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"Name":"Tom"}`), value)

	jc = JsonColumn[User]{}
	value, err = jc.Value()
	assert.Nil(t, err)
	assert.Nil(t, value)
}

func TestJsonColumn_Scan(t *testing.T) {
	testCases := []struct {
		name    string
		src     any
		wantErr error
		wantVal User
	}{
		{
			name:    "nil",
			wantErr: errors.New("ekit：JsonColumn.Scan 不支持 src 类型 <nil>"),
		},
		{
			name:    "string",
			src:     `{"Name":"Tom"}`,
			wantVal: User{Name: "Tom"},
		},
		{
			name: "string pointer",
			src: func() string {
				return `{"Name":"Tom"}`
			}(),
			wantVal: User{Name: "Tom"},
		},
		{
			name:    "bytes",
			src:     []byte(`{"Name":"Tom"}`),
			wantVal: User{Name: "Tom"},
		},
		{
			name: "bytes pointer",
			src: func() *[]byte {
				res := []byte(`{"Name":"Tom"}`)
				return &res
			}(),
			wantVal: User{Name: "Tom"},
		},
		{
			name:    "sql.RawBytes",
			src:     sql.RawBytes(`{"Name":"Tom"}`),
			wantVal: User{Name: "Tom"},
		},
		{
			name: "sql.RawBytes pointer",
			src: func() *sql.RawBytes {
				res := sql.RawBytes(`{"Name":"Tom"}`)
				return &res
			}(),
			wantVal: User{Name: "Tom"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jc := &JsonColumn[User]{}
			err := jc.Scan(tc.src)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, jc.Val)
			assert.True(t, jc.Valid)
		})
	}
}

func TestJsonColumn_ScanTypes(t *testing.T) {
	jcSlice := JsonColumn[[]string]{}
	err := jcSlice.Scan(`["a", "b", "c"]`)
	assert.Nil(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, jcSlice.Val)
	val, err := jcSlice.Value()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`["a","b","c"]`), val)

	jcMap := JsonColumn[map[string]string]{}
	err = jcMap.Scan(`{"a":"a value"}`)
	assert.Nil(t, err)
	val, err = jcMap.Value()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"a":"a value"}`), val)
}

func ExampleJsonColumn_Value() {
	js := JsonColumn[User]{Valid: true, Val: User{Name: "Tom"}}
	value, err := js.Value()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(value.([]byte)))
	//Output:
	//{"Name":"Tom"}
}

func ExampleJsonColumn_Scan() {
	js := JsonColumn[User]{}
	err := js.Scan(`{"Name":"Tom"}`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(js.Val)
	// Output:
	// {Tom}
}

type User struct {
	Name string
}
