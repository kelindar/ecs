// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package builtin

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack"
)

//go:generate genny -pkg=builtin -in=$GOFILE -out=z_components_test.go gen "TType=float32,float64,int16,int32,int64,uint16,uint32,uint64"

func Test_ProviderOfTType(t *testing.T) {
	arr := NewProviderOfTType()
	assert.NotNil(t, arr)
	assert.Equal(t, reflect.TypeOf(pageOfTType{}.data).Elem(), arr.Type())

	index1 := arr.Add(TType(123))
	index2 := arr.Add(TType(123))

	{
		count := 0
		arr.View(func(_ *TType) {
			count++
		})
		assert.Equal(t, 2, count)
	}

	assert.Equal(t, TType(123), arr.ViewAt(index1))
	arr.UpdateAt(index2, func(v *TType) {
		*v = 888
	})
	assert.Equal(t, TType(888), arr.ViewAt(index2))

	arr.RemoveAt(index1)
	arr.RemoveAt(index2)

	{
		count := 0
		arr.Update(func(_ *TType) {
			count++
		})
		assert.Equal(t, 0, count)
	}
}

func Test_CodecOfTType(t *testing.T) {
	original := NewProviderOfTType()
	decoded := NewProviderOfTType()

	// Encode the buffer
	var encoded bytes.Buffer
	err := msgpack.NewEncoder(&encoded).Encode(original)

	// Decode from the buffer
	dec := msgpack.NewDecoder(bytes.NewBuffer(encoded.Bytes()))
	err = dec.Decode(decoded)
	assert.NoError(t, err)
	assert.Equal(t, original, decoded)
}
