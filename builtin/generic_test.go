package builtin

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v4"
)

func Test_Generic(t *testing.T) {
	arr := NewProviderOfTType()
	assert.NotNil(t, arr)

	for i := 0; i < 150; i++ {
		arr.Add("zero")
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

// Benchmark_Component/add-8         	      10	 108778350 ns/op	79380153 B/op	 2000019 allocs/op
// Benchmark_Component/view-8        	     880	   1398878 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Component/update-8      	     897	   1372576 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Component/at-8          	      40	  29211575 ns/op	       0 B/op	       0 allocs/op
func Benchmark_Component(b *testing.B) {
	const size = 1000000

	b.Run("add", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			array := NewProviderOfInt64()
			for i := 0; i < size; i++ {
				array.Add(1)
			}
		}
	})

	array := NewProviderOfInt64()
	for i := 0; i < size; i++ {
		array.Add(1)
	}

	b.Run("view", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			array.View(func(v *int64) {
				return
			})
		}
	})

	b.Run("update", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			x := int64(123)
			array.Update(func(v *int64) {
				*v = x
				return
			})
		}
	})

	b.Run("at", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			for i := 0; i < size; i++ {
				array.ViewAt(i)
			}
		}
	})

}

func Benchmark_Codec(b *testing.B) {
	const size = 1000000
	array := NewProviderOfInt64()
	for i := 0; i < size; i++ {
		array.Add(1)
	}

	b.Run("encode", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			var buf bytes.Buffer
			enc := msgpack.NewEncoder(&buf)
			if err := enc.Encode(array); err != nil {
				b.Fatal(err)
			}
		}
	})

	var encoded bytes.Buffer
	if err := msgpack.NewEncoder(&encoded).Encode(array); err != nil {
		b.Fatal(err)
	}

	b.Run("decode", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			arr := NewProviderOfInt64()
			buf := bytes.NewBuffer(encoded.Bytes())
			dec := msgpack.NewDecoder(buf)
			if err := dec.Decode(arr); err != nil {
				b.Fatal(err)
			}
		}
	})

}
