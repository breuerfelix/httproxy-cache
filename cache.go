package main

type Cache[T any] interface {
	Set(string, *T)
	Get(string) *T
	Delete(string)
}
