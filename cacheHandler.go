package main

import "fmt"

type ChacheHandler interface {
	add(key string, value []byte) error
	get(key string) ([]byte, error)
}

type InMemoryHandler struct {
	data map[string][]byte
}

func NewInMemoryHandler() InMemoryHandler {
	return InMemoryHandler{
		data: make(map[string][]byte),
	}
}

func (i *InMemoryHandler) add(key string, value []byte) error {
	i.data[key] = value
	return nil
}

func (i *InMemoryHandler) get(key string) ([]byte, error) {
	res, ok := i.data[key]
	if !ok {
		return []byte{}, fmt.Errorf("Value for %s not found", key)
	}
	return res, nil

}
