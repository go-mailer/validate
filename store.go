package validate

import (
	"time"
)

// StoreItem 存储的数据项
type StoreItem struct {
	Email        string        // 邮箱
	ValidateCode string        // 验证码
	Time         time.Time     // 存储时间
	Expire       time.Duration // 过期时间
}

// Store 提供线程安全的数据项存储接口，支持自动GC过期的元素
type Store interface {
	// Put 将元素放入存储，返回存储的ID
	// 如果存储发生异常，则返回错误
	Put(item *StoreItem) (int64, error)

	// Take 根据ID取出一个元素(存储内部将移除取出的元素)
	// 如果取出元素发生异常，则返回错误
	Take(id int64) (*StoreItem, error)
}
