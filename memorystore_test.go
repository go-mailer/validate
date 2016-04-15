package validate

import (
	"testing"
	"time"
)

func TestMemoryStoreTakeByID(t *testing.T) {
	store := NewMemoryStore(time.Second * 10)
	item := DataItem{
		Email:      "xxx@gmail.com",
		Code:       "123",
		Expire:     time.Second * 1,
		CreateTime: time.Now(),
	}
	id, err := store.Put(item)
	if err != nil {
		t.Error(err)
		return
	}
	tItem, err := store.TakeByID(id)
	if err != nil {
		t.Error(err)
		return
	}
	if tItem.Email != item.Email {
		t.Error("Get item error:", tItem)
	}
	tItem, err = store.TakeByID(id)
	if err == nil || tItem != nil {
		t.Error("Take item error:", tItem)
	}
}

func TestMemoryStoreTakeByEmailAndCode(t *testing.T) {
	store := NewMemoryStore(time.Second * 10)
	item := DataItem{
		Email:      "xxx@gmail.com",
		Code:       "123",
		Expire:     time.Second * 1,
		CreateTime: time.Now(),
	}
	id, err := store.Put(item)
	if err != nil {
		t.Error(err)
		return
	}
	tItem, err := store.TakeByEmailAndCode("xxx@gmail.com", "123")
	if err != nil {
		t.Error(err)
		return
	}
	if tItem.ID != id {
		t.Error("Get item error:", tItem)
	}
	tItem, err = store.TakeByEmailAndCode("xxx@gmail.com", "123")
	if err == nil || tItem != nil {
		t.Error("Take item error:", tItem)
	}
}

func TestMemoryStoreGC(t *testing.T) {
	store := NewMemoryStore(time.Millisecond * 100)
	item := DataItem{
		Email:      "xxx@gmail.com",
		Code:       "123",
		Expire:     time.Millisecond * 200,
		CreateTime: time.Now(),
	}
	id, err := store.Put(item)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Millisecond * 500)
	tItem, err := store.TakeByID(id)
	if err == nil || tItem != nil {
		t.Error("Take item error:", tItem)
	}
}
