package validate

import (
	"container/list"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// DefaultMemoryGCInterval 默认内存存储GC间隔
	DefaultMemoryGCInterval = time.Second * 60
)

// NewMemoryStore 创建基于内存存储的存储实例
func NewMemoryStore(gcInterval time.Duration) Store {
	if gcInterval == 0 {
		gcInterval = DefaultMemoryGCInterval
	}
	memStore := &MemoryStore{
		idList:     list.New(),
		data:       make(map[int64]*StoreItem),
		gcInterval: gcInterval,
	}
	go memStore.gc()
	return memStore
}

// MemoryStore 提供内存存储
type MemoryStore struct {
	sync.RWMutex
	globalID   int64
	idList     *list.List
	data       map[int64]*StoreItem
	gcInterval time.Duration
}

// gc 执行过期元素检测
func (this *MemoryStore) gc() {
	time.AfterFunc(this.gcInterval, func() {
		defer this.gc()
		for {
			this.RLock()
			ele := this.idList.Front()
			if ele == nil {
				this.RUnlock()
				break
			}
			id := ele.Value.(int64)
			item := this.data[id]
			this.RUnlock()
			if (item.Time.UnixNano() + int64(item.Expire/time.Nanosecond)) < time.Now().UnixNano() {
				this.Lock()
				this.idList.Remove(ele)
				delete(this.data, id)
				this.Unlock()
				continue
			}
			break
		}
	})
}

// Put Put item
func (this *MemoryStore) Put(item *StoreItem) (int64, error) {
	this.Lock()
	defer this.Unlock()
	atomic.AddInt64(&this.globalID, 1)
	this.idList.PushBack(this.globalID)
	this.data[this.globalID] = item
	return this.globalID, nil
}

// Take Take item by ID
func (this *MemoryStore) Take(id int64) (*StoreItem, error) {
	this.RLock()
	var takeEle *list.Element
	for ele := this.idList.Back(); ele != nil; ele = ele.Prev() {
		if ele.Value.(int64) == id {
			takeEle = ele
			break
		}
	}
	if takeEle == nil {
		this.RUnlock()
		return nil, errors.New("Item not found")
	}
	item := this.data[id]
	this.RUnlock()
	this.Lock()
	delete(this.data, id)
	this.idList.Remove(takeEle)
	this.Unlock()
	return item, nil
}
