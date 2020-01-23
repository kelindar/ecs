// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package ecs

import (
	"fmt"
	"reflect"
	"sync"
)

// ComponentType represents a type of the component. This is simply an alias.
type ComponentType = reflect.Type

// Pooler represents that the contract that the component pooler should implement.
type Pooler interface {
	Type() ComponentType
	Add(interface{}) int
	RemoveAt(int)
}

type System interface {
	Name() string
	Priority() int
	Run(*Manager) error
	Stop() error
}

// Manager represents a manager of entities, components and systems.
type Manager struct {
	events
	lock  sync.RWMutex
	pools map[ComponentType]Pooler
	items map[string]map[Serial]*Entity
}

// NewManager returns a new manager instance.
func NewManager() *Manager {
	return &Manager{
		pools: make(map[ComponentType]Pooler),
		items: make(map[string]map[Serial]*Entity),
	}
}

// --------------------------- Manage Entities ----------------------------

// Attach attaches an entity with the set of components.
func (m *Manager) Attach(entity *Entity, components ...interface{}) error {
	for _, part := range components {
		typ := reflect.TypeOf(part)
		pool, ok := m.pools[typ]
		if !ok {
			return fmt.Errorf("type %v is not a valid component", typ.String())
		}

		// Add the part to the pool and keep the index
		entity.parts[typ] = handle{pool, pool.Add(part)}
	}

	// Attach to the registry
	m.lock.Lock()
	if _, ok := m.items[entity.group]; !ok {
		m.items[entity.group] = make(map[Serial]*Entity, 16)
	}
	m.items[entity.group][entity.serial] = entity
	m.lock.Unlock()
	return nil
}

// Detach detaches an entity from the manager and frees the components.
func (m *Manager) Detach(entity *Entity) {
	for _, h := range entity.parts {
		h.mem.RemoveAt(h.idx)
	}

	// Detach from the registry
	m.lock.Lock()
	if group, ok := m.items[entity.group]; ok {
		delete(group, entity.serial)
	}
	m.lock.Unlock()
}

// ---------------------- Manage Component Pools -------------------------

// AddPool registers one or more component pools to the manager.
func (m *Manager) AddPool(pool ...Pooler) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, p := range pool {
		m.pools[p.Type()] = p
	}
}

// RemovePool unregisters one or more component pools from the managers
func (m *Manager) RemovePool(pool ...Pooler) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, p := range pool {
		delete(m.pools, p.Type())
	}
}
