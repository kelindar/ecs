package provider

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v4"
)

func Test_Any(t *testing.T) {
	arr := NewProviderOfAny()
	assert.NotNil(t, arr)

	for i := 0; i < 150; i++ {
		arr.Add("zero")
	}

	count := 0
	arr.View(func(_ *Any) {
		count++
	})

	assert.Equal(t, 3, len(arr.page))
	assert.Equal(t, 150, count)
}

func Test_Page(t *testing.T) {
	var page pageOfAny

	assert.Equal(t, 0, page.Add(nil))
	assert.Equal(t, 1, page.Add(nil))
	assert.Equal(t, 2, page.Add(nil))
	page.Del(1)
	assert.Equal(t, 1, page.Add(nil))
	assert.Equal(t, 3, page.Add(nil))

	count := 0
	page.Range(func(*Any) {
		count++
	})
	assert.Equal(t, 4, count)
	for i := 0; i < 60; i++ {
		page.Add(i)
	}
	assert.Equal(t, true, page.IsFull())
}

// Benchmark_Component/add-8         	      18	  56683200 ns/op	39380796 B/op	      20 allocs/op
// Benchmark_Component/view-8        	     786	   1522678 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Component/update-8      	     786	   1516294 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Component/at-8          	      42	  28329581 ns/op	       0 B/op	       0 allocs/op
func Benchmark_Component(b *testing.B) {
	const size = 1000000
	element := Any(1)

	b.Run("add", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			array := NewProviderOfAny()
			for i := 0; i < size; i++ {
				array.Add(element)
			}
		}
	})

	array := NewProviderOfAny()
	for i := 0; i < size; i++ {
		array.Add(element)
	}

	b.Run("view", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			array.View(func(v *Any) {
				return
			})
		}
	})

	b.Run("update", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			x := int64(123)
			array.Update(func(v *Any) {
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

// Benchmark_Codec/encode-8         	      21	  54188300 ns/op	27158966 B/op	   15647 allocs/op
// Benchmark_Codec/decode-8         	      16	  66944169 ns/op	16656464 B/op	      18 allocs/op
func Benchmark_Codec(b *testing.B) {
	const size = 1000000
	element := Any(1)

	array := NewProviderOfAny()
	for i := 0; i < size; i++ {
		array.Add(element)
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
			arr := NewProviderOfAny()
			buf := bytes.NewBuffer(encoded.Bytes())
			dec := msgpack.NewDecoder(buf)
			if err := dec.Decode(arr); err != nil {
				b.Fatal(err)
			}
		}
	})

}
