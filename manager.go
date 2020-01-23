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

// Provider represents the contract that a component provider should implement.
type Provider interface {
	Type() ComponentType
	Add(interface{}) int
	RemoveAt(int)
}

// System represents the contract that a system should implement.
type System interface {
	Name() string
	Start(*Manager) error
	Close() error
}

// Manager represents a manager of entities, components and systems.
type Manager struct {
	events
	lock  sync.RWMutex
	pools map[ComponentType]Provider
	items map[string]map[Serial]*Entity
	sys   []System
}

// NewManager returns a new manager instance.
func NewManager() *Manager {
	return &Manager{
		pools: make(map[ComponentType]Provider),
		items: make(map[string]map[Serial]*Entity),
		sys:   make([]System, 0, 8),
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

// RegisterProvider registers one or more component pools to the manager.
func (m *Manager) RegisterProvider(providers ...Provider) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, p := range providers {
		m.pools[p.Type()] = p
	}
}

// UnregisterProvider unregisters one or more component pools from the managers
func (m *Manager) UnregisterProvider(providers ...Provider) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, p := range providers {
		delete(m.pools, p.Type())
	}
}

// -------------------------- Manage Systems -----------------------------

// RegisterSystem registers one or more systems to the manager.
func (m *Manager) RegisterSystem(systems ...System) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Start and append all systems
	for _, s := range systems {
		if err := s.Start(m); err != nil {
			return err
		}
		m.sys = append(m.sys, s)
	}
	return nil
}

// UnregisterSystem unregisters one or more systems from the managers.
func (m *Manager) UnregisterSystem(systems ...System) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Filter only valid systems and close the target systems
	var filtered []System
	for _, s := range m.sys {
		found := false
		for _, x := range systems {
			if x == s {
				found = true
				if err := s.Close(); err != nil {
					return err
				}
				break
			}
		}

		if !found {
			filtered = append(filtered, s)
		}
	}
	m.sys = filtered
	return nil
}
