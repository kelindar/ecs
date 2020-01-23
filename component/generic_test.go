// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package component

import (
	"testing"

	"github.com/kelindar/ecs"
	"github.com/stretchr/testify/assert"
)

func Test_Generic(t *testing.T) {
	arr := ForTType()
	assert.NotNil(t, arr)

	for i := 0; i < 150; i++ {
		arr.Add(ecs.NewEntity(), "zero")
	}

	count := 0
	arr.View(func(_ *TType) {
		count++
	})

	assert.Equal(t, 3, len(arr.page))
	assert.Equal(t, 150, count)
}

func Test_Page(t *testing.T) {
	var page pageOfTType

	assert.Equal(t, 0, page.Add(nil))
	assert.Equal(t, 1, page.Add(nil))
	assert.Equal(t, 2, page.Add(nil))
	page.Del(1)
	assert.Equal(t, 1, page.Add(nil))
	assert.Equal(t, 3, page.Add(nil))

	count := 0
	page.Range(func(*TType) {
		count++
	})
	assert.Equal(t, 4, count)
	for i := 0; i < 60; i++ {
		page.Add(i)
	}
	assert.Equal(t, true, page.IsFull())
}

// Benchmark_Array/best-case-8         	    5474	    217540 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Array/array-range-8       	     918	   1321111 ns/op	       0 B/op	       0 allocs/op
func Benchmark_Array(b *testing.B) {
	const size = 1000000

	b.Run("best-case", func(b *testing.B) {
		v := make([]int64, size)
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			for i := 0; i < size; i++ {
				_ = v[i]
			}
		}
	})

	array := ForInt64()
	for i := 0; i < size; i++ {
		array.Add(ecs.NewEntity(), 1)
	}

	b.Run("array-view", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			array.View(func(v *int64) {
				return
			})
		}
	})

}
