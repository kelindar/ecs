// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package builtin

import (
	"bytes"
	"github.com/vmihailenco/msgpack"
	"reflect"
	"testing"

	"github.com/kelindar/ecs"
	"github.com/stretchr/testify/assert"
)

//go:generate genny -pkg=builtin -in=$GOFILE -out=z_components_test.go gen "TType=float32,float64,int16,int32,int64,uint16,uint32,uint64"

func Test_PoolOfTType(t *testing.T) {
	arr := NewPoolOfTType()
	assert.NotNil(t, arr)
	assert.Equal(t, reflect.TypeOf(arr), arr.Type())

	entity1 := ecs.NewEntity()
	entity2 := ecs.NewEntity()

	arr.Add(entity1, 0)
	arr.Add(entity2, 0)

	{
		count := 0
		arr.View(func(_ *TType) {
			count++
		})
		assert.Equal(t, 2, count)
	}

	entity1.Delete()
	entity2.Delete()

	{
		count := 0
		arr.Update(func(_ *TType) {
			count++
		})
		assert.Equal(t, 0, count)
	}
}

func Test_CodecOfTType(t *testing.T) {
	original := NewPoolOfTType()
	decoded := NewPoolOfTType()

	// Encode the buffer
	var encoded bytes.Buffer
	err := msgpack.NewEncoder(&encoded).Encode(original)

	// Decode from the buffer
	dec := msgpack.NewDecoder(bytes.NewBuffer(encoded.Bytes()))
	err = dec.Decode(decoded)
	assert.NoError(t, err)
	assert.Equal(t, original, decoded)
}
