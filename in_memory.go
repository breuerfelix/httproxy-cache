package main

type InMemoryCache[T any] struct {
	db map[string]*T
}

func NewInMemoryCache[T any]() *InMemoryCache[T] {
	return &InMemoryCache[T]{
		db: make(map[string]*T),
	}
}

func (c *InMemoryCache[T]) Set(key string, data *T) {
	c.db[key] = data
}

func (c *InMemoryCache[T]) Get(key string) *T {
	return c.db[key]
}

func (c *InMemoryCache[T]) Delete(key string) {
	delete(c.db, key)
}
