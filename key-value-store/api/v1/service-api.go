package apiv1

import (
	"errors"
	"sync"
)

var ErrorNoSuchKey = errors.New("no such key")


type KeyValueStorekeeper interface {
	Put(string, string) error
	Get(string) (string, error)
	Delete(string) error
}


type KeyValueStore struct {
	store sync.Map
}

func (kvs *KeyValueStore)Put(key string, value string) error {
	kvs.store.Store(key, value)

	return nil
}

func (kvs *KeyValueStore)Get(key string) (string, error) {

	value, ok := kvs.store.Load(key)

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value.(string), nil
}

func (kvs *KeyValueStore)Delete(key string) error {
	kvs.store.Delete(key)

	return nil
}

func NewKeyValueStore()  KeyValueStorekeeper{
	return new(KeyValueStore)
}