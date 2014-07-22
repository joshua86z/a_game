package token

import (
	"fmt"
	"sync"
)

type Adapter struct {
	sync.Mutex
	Map map[string]string
}

func (this *Adapter) Set(key string, value string) error {
	this.Lock()
	this.make()
	this.Map[key] = value
	this.Unlock()
	return nil
}

func (this *Adapter) Get(key string) (string, error) {
	this.Lock()
	this.make()
	str, ok := this.Map[key]
	this.Unlock()
	if !ok {
		return "", fmt.Errorf("can't find this token : %s", key)
	} else {
		return str, nil
	}
}

func (this *Adapter) Delete(key string) error {
	this.Lock()
	this.make()
	delete(this.Map, key)
	this.Unlock()
	return nil
}

func (this *Adapter) make() {
	if this.Map == nil {
		this.Map = make(map[string]string)
	}
}
