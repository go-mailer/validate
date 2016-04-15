# 基于golang实现的邮箱验证

[![GoDoc](https://godoc.org/github.com/go-mailer/validate?status.svg)](https://godoc.org/github.com/go-mailer/validate)

> 支持令牌的安全存储及验证、自动GC过期的验证信息

## 获取

``` bash
$ go get -v github.com/go-mailer/validate
```

## 使用

``` go
package main

import (
	"fmt"
	"time"

	"github.com/go-mailer/validate"
)

func main() {
	// 创建验证信息存储，每10分钟执行一次GC
	store := validate.NewMemoryStore(time.Second * 60 * 10)
	// 创建验证信息管理器，验证信息的过期时间为1小时
	manager := validate.NewManager(store, validate.Config{Expire: time.Second * 60 * 60})
	// 使用邮箱生成验证令牌
	token, err := manager.GenerateToken("xxx@gmail.com")
	if err != nil {
		panic(err)
	}
	fmt.Println("Token:", token)
	// 验证令牌
	isValid, email, err := manager.Validate(token)
	if err != nil {
		panic(err)
	}
	fmt.Println("Valid:", isValid, ",Email:", email)
}
```

## 输出

```
Token: MS45OGZiNWI0ZTU2NmE2NDVhZjRiYmE2ODNmODY3YzllZQ
Valid: true ,Email: xxx@gmail.com
```

## License

	Copyright 2016.All rights reserved.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

