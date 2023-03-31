package runtimex

import (
	"fmt"
	"testing"
)

func TestNewOrderMap(t *testing.T) {
	x := make(map[int]int)
	for i := 0; i < 10; i++ {
		x[i] = i
	}
	for k, v := range x {
		fmt.Println(k, v)
	}

	mm := NewOrderMap()

	for i := 0; i < 10; i++ {
		mm.Insert(i, i)
	}
	println("order")
	mm.Foreach(func(key, val interface{}) error {
		println(key.(int), val.(int))
		return nil
	})
	mm.Delete(1)
	mm.Delete(2)
	println("after")
	mm.Foreach(func(key, val interface{}) error {
		println(key.(int), val.(int))
		return nil
	})
	if mm.Get("10") != nil {
		println(mm.Get(10).(int))
	}
}
