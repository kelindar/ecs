// Copyright (c) 2018 Josh Baker
// Copyright (c) 2020 Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.
// This is a fork of https://github.com/tidwall/rbang adapted for int32 coordinate system.

package atlas

import (
	"fmt"
	"testing"
	"unsafe"
)

// insert:  1,000,000 ops in 592ms, 1,688,010/sec, 592 ns/op
// search:  1,000,000 ops in 272ms, 3,672,805/sec, 272 ns/op
//replace: 1,000,000 ops in 3484ms, 287,050/sec, 3483 ns/op
//delete:  1,000,000 ops in 1505ms, 664,466/sec, 1504 ns/op
func TestGeoIndex(t *testing.T) {
	t.Run("BenchVarious", func(t *testing.T) {
		Tests.TestBenchVarious(t, &Atlas{}, 1000000)
	})
	t.Run("RandomRects", func(t *testing.T) {
		Tests.TestRandomRects(t, &Atlas{}, 10000)
	})
	t.Run("RandomPoints", func(t *testing.T) {
		Tests.TestRandomPoints(t, &Atlas{}, 10000)
	})
	t.Run("ZeroPoints", func(t *testing.T) {
		Tests.TestZeroPoints(t, &Atlas{})
	})
}

func BenchmarkRandomInsert(b *testing.B) {
	Tests.BenchmarkRandomInsert(b, &Atlas{})
}

type Person struct {
	Name    string
	Age     int64
	Pos     [3]float64
	Car     Car
	Friends []Person
}

type Car struct {
	Speed float32
	Model string
}

func TestCast(t *testing.T) {
	p := Person{
		Name: "Roman",
		Age:  12,
		Pos:  [3]float64{1, 2, 3},
		Car: Car{
			Speed: 50,
			Model: "Tesla",
		},
		Friends: []Person{},
	}

	var sizeOfT = unsafe.Sizeof(p)
	fmt.Printf("%d\n", sizeOfT)
	fmt.Printf("%#v\n", p)

	data := (*(*[1<<31 - 1]byte)(unsafe.Pointer(&p)))[:sizeOfT]
	fmt.Printf("%#v\n", data)

	t2 := (*(*Person)(unsafe.Pointer(&data[0])))
	fmt.Printf("%#v\n", t2)

	//assert.Fail(t, "")
}

func Benchmark_Cast(b *testing.B) {
	p := Person{
		Name: "Roman",
		Age:  12,
		Pos:  [3]float64{1, 2, 3},
		Car: Car{
			Speed: 50,
			Model: "Tesla",
		},
		Friends: []Person{},
	}

	var sizeOfT = unsafe.Sizeof(p)
	var data []byte
	b.Run("encode", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			data = (*(*[1<<31 - 1]byte)(unsafe.Pointer(&p)))[:sizeOfT]
		}
	})

	var x Person
	b.Run("decode", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			x = (*(*Person)(unsafe.Pointer(&data[0])))
		}
	})
	fmt.Printf("%#v\n", x)

}
