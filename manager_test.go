// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package ecs

import (
	"reflect"
	"testing"

	"github.com/kelindar/ecs/builtin"
	"github.com/stretchr/testify/assert"
)

func newTestManager() *Manager {
	m := NewManager()
	m.AttachProvider(builtin.NewProviderOfFloat32())
	m.AttachProvider(builtin.NewProviderOfFloat64())
	return m
}

func TestManager(t *testing.T) {
	m := newTestManager()

	// Attach a couple of entities
	e1 := NewEntity("player", 1)
	e2 := NewEntity("player", 2)
	assert.NoError(t, m.AttachEntity(e1, float64(2.0)))
	assert.NoError(t, m.AttachEntity(e2, float64(10.0)))

	m.DetachEntity(e1)

	e, ok := m.GetEntity(Serial(2))
	assert.True(t, ok)
	assert.Equal(t, e2, e)

}

func TestProviders(t *testing.T) {
	m := newTestManager()

	// Attach and get the provider back
	i64 := builtin.NewProviderOfInt64()
	m.AttachProvider(i64)
	_, ok := m.GetProvider(reflect.TypeOf(int64(1.0)))
	assert.True(t, ok)

	// Detach the provider
	m.DetachProvider(i64)
	_, ok = m.GetProvider(reflect.TypeOf(int64(1.0)))
	assert.False(t, ok)

	count := 0
	m.RangeProviders(func(Provider) bool {
		count++
		return true
	})
	assert.Equal(t, 2, count)
}

func TestRangeEntities(t *testing.T) {
	m := newTestManager()

	// Attach a couple of entities
	assert.NoError(t, m.AttachEntity(NewEntity("player", 1), float64(2.0)))
	assert.NoError(t, m.AttachEntity(NewEntity("player", 2), float64(10.0)))
	assert.NoError(t, m.AttachEntity(NewEntity("item", 3), float64(2.0), float32(1.0)))

	count := 0
	m.RangeEntities(func(e *Entity) bool {
		count++
		return true
	})
	assert.Equal(t, 3, count)

	players := 0
	m.RangeEntitiesByGroup("player", func(e *Entity) bool {
		players++
		return true
	})
	assert.Equal(t, 2, players)
}

func TestSystems(t *testing.T) {
	m := newTestManager()
	m.AttachSystem(&testSystem{name: "dummy"})

	// Attach and get the provider back
	s := &testSystem{name: "test"}
	m.AttachSystem(s)
	_, ok := m.GetSystem("test")
	assert.True(t, ok)

	// Detach the provider
	m.DetachSystem(s)
	_, ok = m.GetSystem("test")
	assert.False(t, ok)

	count := 0
	m.RangeSystems(func(System) bool {
		count++
		return true
	})
	assert.Equal(t, 1, count)
}

type testSystem struct {
	name string
}

func (s *testSystem) Name() string {
	return s.name
}

func (s *testSystem) Start(*Manager) error {
	return nil
}

func (s *testSystem) Close() error {
	return nil
}
