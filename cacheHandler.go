package main

import (
	"fmt"
	"strings"
)

type ChacheHandler interface {
	Add(session, key string, value []byte) error
	Get(session, key string) (DataHolder, error)
	Del(session, key string) error
	GetAllValuesForSession(keyPattern string) ([]DataHolder, error)
}

type InMemoryHandler struct {
	data map[string][]byte
}

type DataHolder struct {
	Value   []byte
	Key     string
	Content string
	Session string
}

func getDataHolderByData(key string, value []byte) DataHolder {
	keys := strings.Split(key, "_")
	return DataHolder{
		Value:   value,
		Key:     key,
		Session: keys[0],
		Content: keys[1],
	}
}

func NewInMemoryHandler() InMemoryHandler {
	return InMemoryHandler{
		data: make(map[string][]byte),
	}
}

func (i *InMemoryHandler) concatKey(session, key string) string {
	return fmt.Sprintf("%s_%s", session, key)
}

func (i *InMemoryHandler) Add(session, key string, value []byte) error {
	i.data[i.concatKey(session, key)] = value
	return nil
}

func (i *InMemoryHandler) Del(session, key string) error {
	delete(i.data, i.concatKey(session, key))
	return nil
}

func (i *InMemoryHandler) GetAllValuesForSession(keyPattern string) ([]DataHolder, error) {
	var result []DataHolder
	for key, value := range i.data {
		if strings.HasPrefix(key, keyPattern) {
			result = append(result, getDataHolderByData(key, value))
		}
	}
	return result, nil
}

func (i *InMemoryHandler) Get(session, key string) (DataHolder, error) {
	res, ok := i.data[i.concatKey(session, key)]
	if !ok {
		return DataHolder{}, fmt.Errorf("Value for %s not found", key)
	}
	return getDataHolderByData(i.concatKey(session, key), res), nil

}
