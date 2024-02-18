package main

import (
// "strings"
)

type DB[T any] struct {
	db    map[string]*DB[T]
	value *T
}

type DeepInMemoryCache[T any] struct {
	db map[string]*DB[T]
}

func NewDeepInMemoryCache[T any]() *DeepInMemoryCache[T] {
	return &DeepInMemoryCache[T]{
		db: make(map[string]*DB[T]),
	}
}

func (c *DeepInMemoryCache[T]) Set(key string, data *T) {
	//for _, k := strings.Split(key, "/") {

	//}
	c.db[key] = &DB[T]{value: data}
}

func (c *DeepInMemoryCache[T]) Get(key string) *T {
	val, ok := c.db[key]
	if !ok || val == nil {
		return nil
	}

	return val.value
}

func (c *DeepInMemoryCache[T]) Delete(key string) {
	delete(c.db, key)
}
