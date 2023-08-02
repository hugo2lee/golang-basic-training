/*
 * @Author: hugo lee hugo2lee@gmail.com
 * @Date: 2023-07-31 11:45
 * @LastEditors: hugo lee hugo2lee@gmail.com
 * @LastEditTime: 2023-08-02 18:35
 * @FilePath: /geektime-basic-go/homework/week1/slice.go
 * @Description:
 *
 * Copyright (c) 2023 by hugo, All Rights Reserved.
 */
package week1

import (
	"errors"
)

var ErrIndexOutOfRange = errors.New("下标超出范围")

// 作业：实现切片的删除操作
// 实现删除切片特定下标元素的方法。

// 要求一：能够实现删除操作就可以。
// 要求二：考虑使用比较高性能的实现。
// 要求三：改造为泛型方法
// 要求四：支持缩容，并旦设计缩容机制。

// DeleteAt 删除指定位置的元素
// 如果下标不是合法的下标，返回 ErrIndexOutOfRange
func DeleteAt[T any](s []T, idx int) ([]T, error) {
	if idx >= len(s) {
		return nil, ErrIndexOutOfRange
	}
	return append(s[:idx], s[idx+1:]...), nil
}

func DeleteAtWithCopy[T any](s []T, idx int) ([]T, error) {
	if idx >= len(s) {
		return nil, ErrIndexOutOfRange
	}
	r := make([]T, len(s)-1)
	copy(r, s[:idx])
	copy(r[idx:], s[idx+1:])
	return r, nil
}

func DeleteAtWithGC[T any](s []T, idx int) ([]T, error) {
	if idx >= len(s) {
		return nil, ErrIndexOutOfRange
	}
	copy(s[idx:], s[idx+1:])
	var a T
	s[len(s)-1] = a // 各类型的零值, nil 才会被GC ?
	return s[:len(s)-1], nil
}
