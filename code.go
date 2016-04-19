package validate

import (
	"time"

	"gopkg.in/LyricTian/lib.v2"
)

// NewCodeValidate 创建CodeValidate的实例
// store 验证信息存储方式
// cfg 配置参数(可使用默认参数)
func NewCodeValidate(store Store, cfg ...Config) *CodeValidate {
	config := Config{
		Expire:  DefaultExpire,
		CodeLen: DefaultCodeLen,
	}
	if len(cfg) > 0 {
		v := cfg[0]
		if v.Expire != 0 {
			config.Expire = v.Expire
		}
		if v.CodeLen != 0 {
			config.CodeLen = v.CodeLen
		}
	}
	return &CodeValidate{
		store:  store,
		config: config,
	}
}

// CodeValidate 提供验证码验证
type CodeValidate struct {
	store  Store
	config Config
}

// Generate 根据邮箱生成验证码，同时返回生成的验证码；
// 如果获取失败，则返回错误
func (this *CodeValidate) Generate(email string) (string, error) {
	if email == "" {
		return "", nil
	}
	item := DataItem{
		Email:      email,
		CreateTime: time.Now(),
		Expire:     this.config.Expire,
		Code:       lib.NewRandom(this.config.CodeLen).NumberAndLetter(),
	}
	_, err := this.store.Put(item)
	if err != nil {
		return "", err
	}
	return item.Code, nil
}

// Validate 检查验证信息是否有效
// 如果验证失败，则返回错误
func (this *CodeValidate) Validate(email, code string) (isValid bool, err error) {
	item, err := this.store.TakeByEmailAndCode(email, code)
	if err != nil {
		return
	}
	if (item.CreateTime.UnixNano() + int64(item.Expire/time.Nanosecond)) < time.Now().UnixNano() {
		return
	}
	isValid = true
	return
}
