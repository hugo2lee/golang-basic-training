/*
 * @Author: hugo lee hugo2lee@gmail.com
 * @Date: 2023-07-31 11:45
 * @LastEditors: hugo lee hugo2lee@gmail.com
 * @LastEditTime: 2023-08-09 11:39
 * @FilePath: /geektime-basic-go/homework/week1/slice.go
 * @Description:
 *
 * Copyright (c) 2023 by hugo, All Rights Reserved.
 */
package week1

import (
	"errors"
	"fmt"
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

func DeleteAtNormal[T any](s []T, idx int) ([]T, error) {
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

func DeleteAt[T any](src []T, index int) ([]T, error) {
	length := len(src)
	if index < 0 || index >= length {
		return nil, fmt.Errorf("ekit: %w, 下标超出范围，长度 %d, 下标 %d",
			ErrIndexOutOfRange, length, index)
	}
	for i := index; i+1 < length; i++ {
		src[i] = src[i+1]
	}

	return src[:length-1], nil
}

// Shrink 这是缩容
func Shrink[T any](src []T) []T {
	c, l := cap(src), len(src)
	n, changed := calCapacity(c, l)
	if !changed {
		return src
	}
	s := make([]T, 0, n)
	s = append(s, src...)
	return s
}

func calCapacity(c, l int) (int, bool) {
	// 容量 <=64 缩不缩都无所谓，因为浪费内存也浪费不了多少
	// 你可以考虑调大这个阈值，或者调小这个阈值
	if c <= 64 {
		return c, false
	}
	// 如果容量大于 2048，但是元素不足一半，
	// 降低为 0.625，也就是 5/8
	// 也就是比一半多一点，和正向扩容的 1.25 倍相呼应
	if c > 2048 && (c/l >= 2) {
		factor := 0.625
		return int(float32(c) * float32(factor)), true
	}
	// 如果在 2048 以内，并且元素不足 1/4，那么直接缩减为一半
	if c <= 2048 && (c/l >= 4) {
		return c / 2, true
	}
	// 整个实现的核心是希望在后续少触发扩容的前提下，一次性释放尽可能多的内存
	return c, false
}
