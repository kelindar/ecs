// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntity(t *testing.T) {
	e := NewEntity("player", 1)
	assert.NotNil(t, e)
	assert.Equal(t, Serial(1), e.ID())
	assert.Equal(t, "player", e.Group())
}
