// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package ecs

// Serial represents a unique identifier for an entity.
type Serial uint32

// Handle represents a handle which points to a pool + index combination.
type handle struct {
	mem Pooler // The memory pooler for this handle
	idx int    // The index within the pool
}

// A Entity is simply a set of components with a unique ID attached to it,
// nothing more. It belongs to any amount of Systems, and has a number of
// Components
type Entity struct {
	serial Serial // The identifier of the entity
	group  string // The group name of the entity
	parts  map[ComponentType]handle
}

// NewEntity creates a new entity.
func NewEntity(group string, id Serial) *Entity {
	return &Entity{
		serial: id,
		group:  group,
		parts:  make(map[ComponentType]handle, 8),
	}
}

// ID returns the ID of the entity.
func (e *Entity) ID() Serial {
	return e.serial
}

// Group returns the group name of the entity.
func (e *Entity) Group() string {
	return e.group
}
