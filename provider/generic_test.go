// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package provider

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v4"
)

func Test_ProviderOfAny(t *testing.T) {
	arr := NewProviderOfAny()
	assert.NotNil(t, arr)
	assert.Equal(t, TypeOfAny, arr.Type())

	var v, v2 Any
	index1 := arr.AddAny(v)
	index2 := arr.AddAny(v2)

	{
		count := 0
		arr.View(func(_ *Any) {
			count++
		})
		assert.Equal(t, 2, count)
	}

	assert.Equal(t, v, arr.ViewAt(index1))
	arr.UpdateAt(index2, func(v *Any) {
		*v = v2
	})
	assert.Equal(t, v2, arr.ViewAt(index2))

	arr.RemoveAt(index1)
	arr.RemoveAt(index2)

	{
		count := 0
		arr.Update(func(_ *Any) {
			count++
		})
		assert.Equal(t, 0, count)
	}
}

func Test_CodecOfAny(t *testing.T) {
	original := NewProviderOfAny()
	decoded := NewProviderOfAny()

	// Encode the buffer
	var encoded bytes.Buffer
	err := msgpack.NewEncoder(&encoded).Encode(original)

	// Decode from the buffer
	dec := msgpack.NewDecoder(bytes.NewBuffer(encoded.Bytes()))
	err = dec.Decode(decoded)
	assert.NoError(t, err)
	assert.Equal(t, original, decoded)
}
