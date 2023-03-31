package runtimex

import (
	"container/list"
)

type OrderMap struct {
	key  *list.List
	data map[interface{}]*list.Element
}
type mapEntry struct {
	key   interface{}
	value interface{}
}

func NewOrderMap() *OrderMap {
	return &OrderMap{
		key:  list.New(),
		data: make(map[interface{}]*list.Element),
	}
}
func (o *OrderMap) Insert(key, value interface{}) {
	entry := mapEntry{
		key:   key,
		value: value,
	}
	e := o.key.PushBack(entry)
	o.data[key] = e
}

func (o *OrderMap) Get(key interface{}) interface{} {
	if o.data[key] == nil {
		return nil
	}
	return o.data[key].Value.(mapEntry).value
}

func (o *OrderMap) Foreach(f func(key, val interface{}) error) {
	for e := o.key.Front(); e != nil; e = e.Next() {
		me := e.Value.(mapEntry)
		err := f(me.key, me.value)
		if err != nil {
			break
		}
	}
}
func (o *OrderMap) Delete(key interface{}) {
	e := o.data[key]
	o.key.Remove(e)
	delete(o.data, key)
}
