// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package ecs

import (
	"context"
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
	Start(context.Context, *Manager) error
	Close() error
}

// Manager represents a manager of entities, components and systems.
type Manager struct {
	lock  sync.RWMutex
	sys   map[string]System
	pools map[ComponentType]Provider
	byids map[Serial]*Entity
	bygrp map[string]map[Serial]*Entity
}

// NewManager returns a new manager instance.
func NewManager() *Manager {
	return &Manager{
		pools: make(map[ComponentType]Provider, 100),
		bygrp: make(map[string]map[Serial]*Entity, 100),
		byids: make(map[Serial]*Entity, 1000000),
		sys:   make(map[string]System, 16),
	}
}

// --------------------------- Manage Entities ----------------------------

// AttachEntity attaches an entity with the set of components.
func (m *Manager) AttachEntity(entity *Entity, components ...interface{}) error {
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
	if _, ok := m.bygrp[entity.group]; !ok {
		m.bygrp[entity.group] = make(map[Serial]*Entity, 16)
	}
	m.bygrp[entity.group][entity.serial] = entity
	m.byids[entity.serial] = entity
	m.lock.Unlock()
	return nil
}

// DetachEntity detaches an entity from the manager and frees the components.
func (m *Manager) DetachEntity(entity *Entity) {
	for _, h := range entity.parts {
		h.mem.RemoveAt(h.idx)
	}

	// Detach from the registry
	m.lock.Lock()
	if group, ok := m.bygrp[entity.group]; ok {
		delete(group, entity.serial)
	}
	delete(m.byids, entity.serial)
	m.lock.Unlock()
}

// RangeEntitiesByGroup iterates over the entities of a specific group.
func (m *Manager) RangeEntitiesByGroup(group string, f func(*Entity) bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if group, ok := m.bygrp[group]; ok {
		for _, e := range group {
			if !f(e) {
				return
			}
		}
	}
}

// RangeEntities iterates over all entities present. Note that this can be slow
// and acquires a read lock, prefer iterating for a single group instead.
func (m *Manager) RangeEntities(f func(*Entity) bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, e := range m.byids {
		if !f(e) {
			return
		}
	}
}

// GetEntity returns the entity by its Serial.
func (m *Manager) GetEntity(id Serial) *Entity {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if e, ok := m.byids[id]; ok {
		return e
	}
	return nil
}

// ---------------------- Manage Component Pools -------------------------

// AttachProvider registers one or more component pools to the manager.
func (m *Manager) AttachProvider(providers ...Provider) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, p := range providers {
		m.pools[p.Type()] = p
	}
}

// DetachProvider unregisters one or more component pools from the managers
func (m *Manager) DetachProvider(providers ...Provider) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, p := range providers {
		delete(m.pools, p.Type())
	}
}

// RangeProviders iterates over all registered component providers.
func (m *Manager) RangeProviders(f func(Provider) bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, p := range m.pools {
		if !f(p) {
			return
		}
	}
}

// GetProvider returns the provider for a specific component type.
func (m *Manager) GetProvider(typ ComponentType) Provider {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if p, ok := m.pools[typ]; ok {
		return p
	}
	return nil
}

// -------------------------- Manage Systems -----------------------------

// AttachSystem registers one or more systems to the manager.
func (m *Manager) AttachSystem(ctx context.Context, systems ...System) error {
	for _, s := range systems {
		if err := s.Start(ctx, m); err != nil {
			return err
		}

		m.lock.Lock()
		m.sys[s.Name()] = s
		m.lock.Unlock()
	}
	return nil
}

// DetachSystem DetachSystem one or more systems from the managers.
func (m *Manager) DetachSystem(systems ...System) error {
	for _, x := range systems {
		if sys, ok := m.sys[x.Name()]; ok {
			m.lock.Lock()
			delete(m.sys, x.Name())
			m.lock.Unlock()

			if err := sys.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

// RangeSystems iterates over all registered systems.
func (m *Manager) RangeSystems(f func(System) bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, s := range m.sys {
		if !f(s) {
			return
		}
	}
}

// GetSystem returns the system by its name.
func (m *Manager) GetSystem(name string) System {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if s, ok := m.sys[name]; ok {
		return s
	}
	return nil
}
