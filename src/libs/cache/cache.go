package cache

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"
	"time"
)

const maxGroup = 10000

type item struct {
	value  interface{}
	expire int64
}

type group struct {
	sync.Mutex
	itemMap map[string]*item
}

type instance struct {
	group [maxGroup]group
}

func (this *instance) GC() {
	for {
		for index := range this.group {

			this.group[index].Lock()
			if this.group[index].itemMap == nil {
				this.group[index].Unlock()
				continue
			}
			for key, item := range this.group[index].itemMap {
				if time.Now().Unix() > item.expire {
					delete(this.group[index].itemMap, key)
				}
			}

			this.group[index].Unlock()
			time.Sleep(time.Second * 60)
		}
	}
}

var hashMap *instance

func Instance() *instance {

	if hashMap == nil {
		hashMap = &instance{}
		go hashMap.GC()
	}

	return hashMap
}

func (this *instance) Get(key string) (interface{}, error) {

	index := find(key)

	this.group[index].Lock()

	if this.group[index].itemMap == nil {

		this.group[index].itemMap = make(map[string]*item)
		this.group[index].Unlock()

		return nil, fmt.Errorf("Can't find the key : '%s'", key)
	}

	var err error
	var result interface{}

	if val, ok := this.group[index].itemMap[key]; !ok {

		err = fmt.Errorf("Can't find the key : '%s'", key)
	} else if time.Now().Unix() >= val.expire {

		delete(this.group[index].itemMap, key)

		err = fmt.Errorf("The key has expired")
	} else {

		result = val.value
	}

	this.group[index].Unlock()

	return result, err
}

func (this *instance) Set(key string, value interface{}, expire int) {

	index := find(key)

	this.group[index].Lock()

	if this.group[index].itemMap == nil {
		this.group[index].itemMap = make(map[string]*item)
	}

	if this.group[index].itemMap[key] == nil {
		this.group[index].itemMap[key] = &item{}
	}

	this.group[index].itemMap[key].value = value
	this.group[index].itemMap[key].expire = time.Now().Add(time.Second * time.Duration(expire)).Unix()

	this.group[index].Unlock()
}

func (this *instance) Delete(key string) {

	index := find(key)

	this.group[index].Lock()

	if this.group[index].itemMap == nil {

		this.group[index].Unlock()
		return
	}

	delete(this.group[index].itemMap, key)

	this.group[index].Unlock()
}

func find(key string) int {
	return int(Hash32([]byte(key)) % maxGroup)
}

func Hash32(key []byte) uint32 {

	length := len(key)
	if length == 0 {
		return 0
	}

	var c1, c2, h, k uint32
	c1 = 0xcc9e2d51
	c2 = 0x1b873593

	buf := bytes.NewBuffer(key)

	nblocks := length / 4
	for i := 0; i < nblocks; i++ {
		binary.Read(buf, binary.LittleEndian, &k)
		k *= c1
		k = (k << 15) | (k >> (32 - 15))
		k *= c2
		h ^= k
		h = (h << 13) | (h >> (32 - 13))
		h = (h * 5) + 0xe6546b64
	}

	k = 0
	tailIndex := nblocks * 4

	switch length & 3 {
	case 3:
		k ^= uint32(key[tailIndex+2]) << 16
		fallthrough
	case 2:
		k ^= uint32(key[tailIndex+1]) << 8
		fallthrough
	case 1:
		k ^= uint32(key[tailIndex])
		k *= c1
		k = (k << 15) | (k >> (32 - 15))
		k *= c2
		h ^= k
	}

	h ^= uint32(length)
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16

	return h
}
