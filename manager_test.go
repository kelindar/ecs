// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package ecs

import (
	"testing"

	"github.com/kelindar/ecs/builtin"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	f64 := builtin.NewProviderOfFloat64()

	m := NewManager()
	m.AttachProvider(f64)
	defer m.DetachProvider(f64)

	e := NewEntity("player", 1)

	err := m.AttachEntity(e, float64(2.0))
	assert.NoError(t, err)

	m.DetachEntity(e)

}
