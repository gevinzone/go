package misc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func TestReturnValue(t *testing.T) {
	// 值相同，但由于值传递有复制，故内存地址不同
	res, p := returnValue()
	fmt.Printf("p: %p\n", &res)
	assert.NotEqual(t, p, unsafe.Pointer(&res))
}

func TestReturnPointer(t *testing.T) {
	// 指针指向的值相同，但由于指针本身经过了复制，故指向指针的指针（二级指针）值不同
	res, p := returnPointer()
	fmt.Printf("p: %p, *p: %p\n", res, &res)
	assert.Equal(t, p, unsafe.Pointer(res))
}

func TestPassValue(t *testing.T) {
	d := demo{id: 1}
	passValue(d)
	fmt.Printf("p: %p\n", &d)
}

func TestPassPointer(t *testing.T) {
	d := &demo{id: 1}
	passPointer(d)
	fmt.Printf("p: %p, *p: %p\n", d, &d)
}

type demo struct {
	id int64
}

func returnValue() (demo, unsafe.Pointer) {
	res := demo{id: 1}
	fmt.Printf("p: %p\n", &res)
	return res, unsafe.Pointer(&res)
}

func returnPointer() (*demo, unsafe.Pointer) {
	res := &demo{id: 1}
	fmt.Printf("p: %p, *p: %p\n", res, &res)
	return res, unsafe.Pointer(res)
}

func passValue(d demo) {
	fmt.Printf("p: %p\n", &d)
}

func passPointer(d *demo) {
	fmt.Printf("p: %p, *p: %p\n", d, &d)
}
