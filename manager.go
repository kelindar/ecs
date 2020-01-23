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

// Manager represents a manager of entities, components and systems.
type Manager struct {
	lock     sync.RWMutex
	pool     map[reflect.Type]Pooler
	entities map[Serial]*Entity
}

// NewManager returns a new manager instance.
func NewManager() *Manager {
	return &Manager{
		pool:     make(map[ComponentType]Pooler),
		entities: make(map[Serial]*Entity),
	}
}

// ---------------------- Manage Component Pools -------------------------

// Register registers one or more component pools to the manager.
func (m *Manager) Register(pool ...Pooler) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, p := range pool {
		m.pool[p.Type()] = p
	}
}

// Unregister unregisters one or more component pools from the managers
func (m *Manager) Unregister(pool ...Pooler) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, p := range pool {
		delete(m.pool, p.Type())
	}
}

// --------------------------- Manage Entities ----------------------------

// Attach attaches an entity with the set of components.
func (m *Manager) Attach(entity *Entity, components ...interface{}) error {
	for _, part := range components {
		typ := reflect.TypeOf(part)
		pool, ok := m.pool[typ]
		if !ok {
			return fmt.Errorf("type %v is not a valid component", typ.String())
		}

		// Add the part to the pool and keep the index
		entity.parts[typ] = handle{pool, pool.Add(part)}
	}

	return nil
}

// Detach detaches an entity from the manager and frees the components.
func (m *Manager) Detach(entity *Entity) {
	for _, h := range entity.parts {
		h.mem.RemoveAt(h.idx)
	}
}
