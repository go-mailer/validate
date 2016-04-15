package validate

import (
	"testing"
	"time"
)

func TestMemoryStore(t *testing.T) {
	store := NewMemoryStore(time.Second * 10)
	item := &StoreItem{
		Email:        "xxx@gmail.com",
		ValidateCode: "123",
		Expire:       time.Second * 1,
		Time:         time.Now(),
	}
	id, err := store.Put(item)
	if err != nil {
		t.Error(err)
		return
	}
	tItem, err := store.Take(id)
	if err != nil {
		t.Error(err)
		return
	}
	if tItem.Email != item.Email {
		t.Error("Get item error:", tItem)
	}
	tItem, err = store.Take(id)
	if err == nil || tItem != nil {
		t.Error("Take item error:", tItem)
	}
}

func TestMemoryStoreGC(t *testing.T) {
	store := NewMemoryStore(time.Millisecond * 100)
	item := &StoreItem{
		Email:        "xxx@gmail.com",
		ValidateCode: "123",
		Expire:       time.Millisecond * 200,
		Time:         time.Now(),
	}
	id, err := store.Put(item)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Millisecond * 500)
	tItem, err := store.Take(id)
	if err == nil || tItem != nil {
		t.Error("Take item error:", tItem)
	}
}
