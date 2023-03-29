package gocache

import (
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func TestContainer(t *testing.T) {
	parallel := 50
	container := InitContainer(parallel / 2)
	callback := make(chan error, parallel)
	for i := 0; i < parallel; i++ {
		go testContainerInParallel(container, callback)
	}
	for i := 0; i < parallel; i++ {
		res := <-callback
		if res != nil {
			t.Fatalf("test container in parallel error, %s", res)
		}
	}
}

func testContainerInParallel(container Container, callback chan error) {
	key := uuid.New().String()
	valueSet := uuid.New().String()
	container.Set(key, valueSet)
	valueGet := container.Get(key)
	if valueGet != valueSet {
		callback <- fmt.Errorf("wrong get/set")
		return
	}
	container.Del(key)
	valueDel := container.Get(key)
	if valueDel != "" {
		callback <- fmt.Errorf("wrong get/del")
		return
	}
	callback <- nil
}
