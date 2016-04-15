package validate

import "time"

const (
	// DefaultExpire 默认过期时间（2个小时）
	DefaultExpire = time.Hour * 2
	// DefaultRandomCodeLen 默认随机码的长度
	DefaultRandomCodeLen = 6
)

// Config 邮箱验证的配置参数
type Config struct {
	Expire        time.Duration // 过期时间
	RandomCodeLen int           // 随机码的长度
}
