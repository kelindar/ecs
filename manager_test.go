// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package ecs

import (
	"testing"

	"github.com/kelindar/ecs/builtin"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	f64 := builtin.NewPoolOfFloat64()

	m := NewManager()
	m.AddPool(f64)
	defer m.RemovePool(f64)

	e := NewEntity("player", 1)

	err := m.Attach(e, float64(2.0))
	assert.NoError(t, err)

	m.Detach(e)
}
