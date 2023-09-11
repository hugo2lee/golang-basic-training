/*
 * @Author: hugo2lee hugo2lee@gmail.com
 * @Date: 2023-09-03 13:46
 * @LastEditors: hugo2lee hugo2lee@gmail.com
 * @LastEditTime: 2023-09-11 10:16
 * @FilePath: /geektime-basic-go/webook/internal/repository/cache/local_code.go
 * @Description:
 *
 * Copyright (c) 2023 by hugo, All Rights Reserved.
 */
package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type CodeCache interface {
	Set(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

// LocalCodeCache 基于 local 的实现
type LocalCodeCache struct {
	cache *cache.Cache
	lock  sync.Mutex // go-cache里面有锁, 业务还是要加锁
}

func NewLocalCodeCache(ca *cache.Cache) CodeCache {
	return &LocalCodeCache{
		cache: ca,
	}
}

// Set 如果该手机在该业务场景下，验证码不存在（都已经过期），那么发送
// 如果已经有一个验证码，但是发出去已经一分钟了，允许重发
// 如果已经有一个验证码，但是没有过期时间，说明有不知名错误
// 如果已经有一个验证码，但是发出去不到一分钟，不允许重发
// 验证码有效期 10 分钟
func (c *LocalCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ttl, ok := c.cache.GetWithExpiration(c.key(biz, phone))

	if !ok || time.Since(ttl) > 60*time.Second {
		if err := c.cache.Add(c.key(biz, phone), code, time.Second*600); err != nil {
			return ErrUnknownForCode
		}
		if err := c.cache.Add(c.Cntkey(biz, phone), 3, time.Second*600); err != nil {
			return ErrUnknownForCode
		}
		return nil
	} else if ok && ttl.Equal(time.Time{}) {
		return ErrUnknownForCode
	} else {
		return ErrCodeSendTooMany
	}
}

// Verify 验证验证码
// 如果验证码是一致的，那么删除
// 如果验证码不一致，那么保留的
func (c *LocalCodeCache) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	cnt, _, cntOk := c.cache.GetWithExpiration(c.Cntkey(biz, phone))
	if !cntOk {
		return false, ErrUnknownForCode
	}
	cn, ok := cnt.(int)
	if !ok {
		return false, ErrUnknownForCode
	}
	if cn <= 0 {
		return false, ErrCodeVerifyTooManyTimes
	}

	expectedCode, _, keyOk := c.cache.GetWithExpiration(c.key(biz, phone))
	if !keyOk {
		return false, ErrUnknownForCode
	}
	expected, ok := expectedCode.(string)
	if !ok {
		return false, ErrUnknownForCode
	}

	if expected == inputCode {
		c.cache.Set(c.Cntkey(biz, phone), -1, 0)
		return true, nil
	} else {
		c.cache.Set(c.Cntkey(biz, phone), cn-1, 0)
		return false, nil
	}
}

func (c *LocalCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (c *LocalCodeCache) Cntkey(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s:cnt", biz, phone)
}
