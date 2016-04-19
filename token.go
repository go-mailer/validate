package validate

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/LyricTian/lib.v2"
)

// NewTokenValidate 创建TokenValidate的实例
// store 验证信息存储方式
// cfg 配置参数(可使用默认参数)
func NewTokenValidate(store Store, cfg ...Config) *TokenValidate {
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
	return &TokenValidate{
		store:  store,
		config: config,
	}
}

// TokenValidate 提供令牌验证
type TokenValidate struct {
	store  Store
	config Config
}

// getToken 生成token
func (this *TokenValidate) getToken(item *DataItem) (string, error) {
	var itemBuf bytes.Buffer
	itemBuf.WriteString(item.Email)
	itemBuf.WriteByte('\n')
	itemBuf.WriteString(item.Code)
	itemBuf.WriteByte('\n')
	itemBuf.WriteString(strconv.FormatInt(item.CreateTime.Unix(), 10))
	token, err := lib.NewEncryption(itemBuf.Bytes()).MD5()
	if err != nil {
		return "", err
	}
	itemBuf.Reset()
	return token, nil
}

// Generate 根据邮箱生成令牌，同时返回生成的令牌；
// 如果获取失败，则返回错误
func (this *TokenValidate) Generate(email string) (string, error) {
	if email == "" {
		return "", nil
	}
	item := DataItem{
		Email:      email,
		CreateTime: time.Now(),
		Expire:     this.config.Expire,
		Code:       lib.NewRandom(this.config.CodeLen).NumberAndLetter(),
	}
	id, err := this.store.Put(item)
	if err != nil {
		return "", err
	}
	token, err := this.getToken(&item)
	if err != nil {
		return "", err
	}
	val := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d.%s", id, token)))
	return strings.TrimRight(val, "="), nil
}

// parseToken 解析token
func (this *TokenValidate) parseToken(token string) (id int64, encryptVal string, err error) {
	token = token + strings.Repeat("=", 4-len(token)%4)
	v, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return
	}
	tokenVal := strings.SplitN(string(v), ".", 2)
	if len(tokenVal) != 2 {
		err = errors.New("Token is invalid")
		return
	}
	id, err = strconv.ParseInt(tokenVal[0], 10, 64)
	if err != nil {
		return
	}
	encryptVal = tokenVal[1]
	return
}

// Validate 验证令牌是否有效，同时返回有效的邮箱地址;
// 如果验证失败，则返回错误
func (this *TokenValidate) Validate(token string) (isValid bool, email string, err error) {
	id, encryptVal, err := this.parseToken(token)
	if err != nil {
		return
	}
	item, err := this.store.TakeByID(id)
	if err != nil {
		return
	}
	idToken, err := this.getToken(item)
	if err != nil {
		return
	}
	if encryptVal != idToken ||
		(item.CreateTime.UnixNano()+int64(item.Expire/time.Nanosecond)) < time.Now().UnixNano() {
		return
	}
	isValid = true
	email = item.Email
	return
}
