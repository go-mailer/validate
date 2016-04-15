package validate

import (
	"container/list"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// NewMemoryStore 创建基于内存存储的存储实例
func NewMemoryStore(gcInterval time.Duration) Store {
	if gcInterval == 0 {
		gcInterval = DefaultGCInterval
	}
	memStore := &MemoryStore{
		data:       list.New(),
		gcInterval: gcInterval,
	}
	go memStore.gc()
	return memStore
}

// MemoryStore 提供内存存储
type MemoryStore struct {
	sync.RWMutex
	globalID   int64
	gcInterval time.Duration
	data       *list.List
}

func (this *MemoryStore) gc() {
	time.AfterFunc(this.gcInterval, func() {
		defer this.gc()
		for {
			this.RLock()
			ele := this.data.Front()
			if ele == nil {
				this.RUnlock()
				break
			}
			item := ele.Value.(DataItem)
			this.RUnlock()
			if (item.CreateTime.UnixNano() + int64(item.Expire/time.Nanosecond)) < time.Now().UnixNano() {
				this.Lock()
				this.data.Remove(ele)
				this.Unlock()
				continue
			}
			break
		}
	})
}

// Put Put item
func (this *MemoryStore) Put(item DataItem) (int64, error) {
	this.Lock()
	defer this.Unlock()
	atomic.AddInt64(&this.globalID, 1)
	item.ID = this.globalID
	this.data.PushBack(item)
	return item.ID, nil
}

// Take Take item by ID
func (this *MemoryStore) TakeByID(id int64) (*DataItem, error) {
	this.RLock()
	var takeEle *list.Element
	for ele := this.data.Back(); ele != nil; ele = ele.Prev() {
		item := ele.Value.(DataItem)
		if item.ID == id {
			takeEle = ele
			break
		}
	}
	if takeEle == nil {
		this.RUnlock()
		return nil, errors.New("Item not found")
	}
	item := takeEle.Value.(DataItem)
	this.RUnlock()
	this.Lock()
	this.data.Remove(takeEle)
	this.Unlock()
	return &item, nil
}

// Take Take item by email and code
func (this *MemoryStore) TakeByEmailAndCode(email, code string) (*DataItem, error) {
	this.RLock()
	var takeEle *list.Element
	for ele := this.data.Back(); ele != nil; ele = ele.Prev() {
		item := ele.Value.(DataItem)
		if item.Email == email && item.Code == code {
			takeEle = ele
			break
		}
	}
	if takeEle == nil {
		this.RUnlock()
		return nil, errors.New("Item not found")
	}
	item := takeEle.Value.(DataItem)
	this.RUnlock()
	this.Lock()
	this.data.Remove(takeEle)
	this.Unlock()
	return &item, nil
}
