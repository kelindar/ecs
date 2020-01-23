// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package component

import (
	"testing"

	"github.com/kelindar/ecs"
	"github.com/stretchr/testify/assert"
)

//go:generate genny -in=$GOFILE -out=z_components_test.go gen "TType=float32,float64,int16,int32,int64,uint16,uint32,uint64"

func Test_TType(t *testing.T) {
	arr := ForTType()
	assert.NotNil(t, arr)

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
