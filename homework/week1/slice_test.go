/*
 * @Author: hugo lee hugo2lee@gmail.com
 * @Date: 2023-07-31 11:45
 * @LastEditors: hugo lee hugo2lee@gmail.com
 * @LastEditTime: 2023-08-02 18:33
 * @FilePath: /geektime-basic-go/homework/week1/slice_test.go
 * @Description:
 *
 * Copyright (c) 2023 by hugo, All Rights Reserved.
 */
package week1

import (
	"reflect"
	"testing"
)

func TestDeleteAt(t *testing.T) {
	type sample[T any] struct {
		name    string
		s       []T
		idx     int
		want    []T
		wantErr bool
	}
	tests := []sample[int]{
		{
			"delete first",
			[]int{1, 2, 3},
			0,
			[]int{2, 3},
			false,
		},
		{
			"delete middle",
			[]int{1, 2, 3},
			1,
			[]int{1, 3},
			false,
		},
		{
			"delete last",
			[]int{1, 2, 3},
			2,
			[]int{1, 2},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteAt(tt.s, tt.idx)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteAtWithCopy(t *testing.T) {
	type sample[T any] struct {
		name    string
		s       []T
		idx     int
		want    []T
		wantErr bool
	}
	tests := []sample[int]{
		{
			"delete first",
			[]int{1, 2, 3},
			0,
			[]int{2, 3},
			false,
		},
		{
			"delete middle",
			[]int{1, 2, 3},
			1,
			[]int{1, 3},
			false,
		},
		{
			"delete last",
			[]int{1, 2, 3},
			2,
			[]int{1, 2},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteAtWithCopy(tt.s, tt.idx)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteAtWithCopy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteAtWithCopy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteAtWithGC(t *testing.T) {
	type sample[T any] struct {
		name    string
		s       []T
		idx     int
		want    []T
		wantErr bool
	}
	tests := []sample[int]{
		{
			"delete first",
			[]int{1, 2, 3},
			0,
			[]int{2, 3},
			false,
		},
		{
			"delete middle",
			[]int{1, 2, 3},
			1,
			[]int{1, 3},
			false,
		},
		{
			"delete last",
			[]int{1, 2, 3},
			2,
			[]int{1, 2},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteAtWithGC(tt.s, tt.idx)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteAtWithGC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteAtWithGC() = %v, want %v", got, tt.want)
			}
		})
	}
}
