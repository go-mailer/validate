package validate

// Store 提供线程安全的验证信息存储接口，支持自动GC过期的元素
type Store interface {
	// Put 将元素放入存储，返回存储的ID
	// 如果存储发生异常，则返回错误
	Put(item DataItem) (int64, error)

	// TakeByID 根据ID取出一个元素，返回取出的元素
	// 如果取出元素发生异常，则返回错误
	TakeByID(id int64) (*DataItem, error)

	// TakeByEmailAndCode 根据验证码和验证码取出一个元素，返回取出的元素
	// 如果取出元素发生异常，则返回错误
	TakeByEmailAndCode(email, code string) (*DataItem, error)
}
