// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package ecs

import (
	"context"
	"testing"

	"github.com/kelindar/ecs/provider"
	"github.com/stretchr/testify/assert"
)

func newTestManager() *Manager {
	m := NewManager()
	m.AttachProvider(provider.NewProviderOfPoint())
	return m
}

func TestManager(t *testing.T) {
	m := newTestManager()

	// Attach a couple of entities
	e1 := NewEntity("player", 1)
	e2 := NewEntity("player", 2)
	assert.NoError(t, m.AttachEntity(e1, provider.Point{X: 1}))
	assert.NoError(t, m.AttachEntity(e2, provider.Point{X: 2}))

	m.DetachEntity(e1)

	e := m.GetEntity(Serial(2))
	assert.NotNil(t, e)
	assert.Equal(t, e2, e)

}

func TestProviders(t *testing.T) {
	m := NewManager()

	// Attach and get the provider back
	p := provider.NewProviderOfPoint()
	m.AttachProvider(p)
	out := m.GetProvider(provider.TypeOfPoint)
	assert.NotNil(t, out)

	count := 0
	m.RangeProviders(func(Provider) bool {
		count++
		return true
	})
	assert.Equal(t, 1, count)

	// Detach the provider
	m.DetachProvider(p)
	out = m.GetProvider(provider.TypeOfPoint)
	assert.Nil(t, out)

	count = 0
	m.RangeProviders(func(Provider) bool {
		count++
		return true
	})
	assert.Equal(t, 0, count)
}

func TestRangeEntities(t *testing.T) {
	m := newTestManager()

	// Attach a couple of entities
	assert.NoError(t, m.AttachEntity(NewEntity("player", 1), provider.Point{X: 1}))
	assert.NoError(t, m.AttachEntity(NewEntity("player", 2), provider.Point{X: 1}))
	assert.NoError(t, m.AttachEntity(NewEntity("item", 3), provider.Point{X: 1}))

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
	m.AttachSystem(context.Background(), &testSystem{name: "dummy"})

	// Attach and get the provider back
	s := &testSystem{name: "test"}
	m.AttachSystem(context.Background(), s)
	out := m.GetSystem("test")
	assert.NotNil(t, out)

	// Detach the provider
	m.DetachSystem(s)
	out = m.GetSystem("test")
	assert.Nil(t, out)

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

func (s *testSystem) Start(context.Context, *Manager) error {
	return nil
}

func (s *testSystem) Close() error {
	return nil
}
