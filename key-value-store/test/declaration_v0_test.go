package test

import (
	"errors"
	api "github.com/lissdx/go-examples/key-value-store/pkg/api/v0"
	"testing"
)

func TestPut(t *testing.T) {
	const key = "create-key"
	const value = "create-value"

	var val interface{}
	var contains bool

	defer delete(api.Store, key)

	// Sanity check
	_, contains = api.Store[key]
	if contains {
		t.Error("key/value already exists")
	}

	// err should be nil
	err := api.Put(key, value)
	if err != nil {
		t.Error(err)
	}

	val, contains = api.Store[key]
	if !contains {
		t.Error("create failed")
	}

	if val != value {
		t.Error("val/value mismatch")
	}
}

func TestGet(t *testing.T) {
	const key = "read-key"
	const value = "read-value"

	var val interface{}
	var err error

	defer delete(api.Store, key)

	// Read a non-thing
	val, err = api.Get(key)
	if err == nil {
		t.Error("expected an error")
	}
	if !errors.Is(err, api.ErrorNoSuchKey) {
		t.Error("unexpected error:", err)
	}

	api.Store[key] = value

	val, err = api.Get(key)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if val != value {
		t.Error("val/value mismatch")
	}
}

func TestDelete(t *testing.T) {
	const key = "delete-key"
	const value = "delete-value"

	var contains bool

	defer delete(api.Store, key)

	api.Store[key] = value

	_, contains = api.Store[key]
	if !contains {
		t.Error("key/value doesn't exist")
	}

	api.Delete(key)

	_, contains = api.Store[key]
	if contains {
		t.Error("Delete failed")
	}
}
