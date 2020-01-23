// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package ecs

import "sync/atomic"

// TODO: when loaded, we need to set the last serial
var lastSerial uint32

// Serial represents a unique identifier for an entity.
type Serial uint32

// next generates next available serial
func next() Serial {
	return Serial(atomic.AddUint32(&lastSerial, 1))
}

// ---------------------------------------------------------------------------------

// A Entity is simply a set of components with a unique ID attached to it,
// nothing more. It belongs to any amount of Systems, and has a number of
// Components
type Entity struct {
	ID     Serial   // The identifier of an entity
	detach []func() // The closer functions that remove the entity from components
}

// NewEntity creates a new Entity with a new unique identifier. It is safe for
// concurrent use.
func NewEntity() *Entity {
	return &Entity{ID: next()}
}

// Attach attaches a dispose function to the entity
func (e *Entity) Attach(remove func()) {
	e.detach = append(e.detach, remove)
}

// Delete deletes the entity and removes it from the components
func (e *Entity) Delete() {
	for _, remove := range e.detach {
		remove()
	}
}
