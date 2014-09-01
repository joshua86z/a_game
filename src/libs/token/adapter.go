package token

import (
	"fmt"
	"sync"
)

func NewAdapter() *Adapter {
	return &Adapter{Map: make(map[string]string)}
}

type Adapter struct {
	sync.RWMutex
	Map map[string]string
}

func (this *Adapter) Set(key string, value string) error {
	this.Lock()
	this.Map[key] = value
	this.Unlock()
	return nil
}

func (this *Adapter) Get(key string) (string, error) {
	this.RLock()
	str, ok := this.Map[key]
	this.RUnlock()
	if !ok {
		return "", fmt.Errorf("can't find this token : %s", key)
	} else {
		return str, nil
	}
}

func (this *Adapter) Delete(key string) error {
	this.Lock()
	delete(this.Map, key)
	this.Unlock()
	return nil
}
