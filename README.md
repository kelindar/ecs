# Example of Entity Component System in Go

## Introduction

This is my attempt to build a high-performance ECS (Entity, Component, System) in pure Go. Having investigated other ECS systems which have been written to date, I realised that none of them is really a pure ECS as the data layout is often using pointers or interfaces, which would discard most of the benefits of such a system.

Instead of being an ECS framework, this repository contains an example of using [kelindar/column](https://github.com/kelindar/column) as the underlying storage of components and [kelindar/tile](https://github.com/kelindar/tile) as a 2D grid engine for spatial queries.

## Entities

The `entities` directory contains various **entities** of the game, things such as players, items, monsters etc. An entity represents something that has a set of **components** (i.e. columns) and is expressed as a **view** over a row that is lazily evaluated by various **systems**.

For example, a mobile entity takes a pointer over a specific row and allows us to read/write columns by simply accessing the properties instead of dealing with the internals of the data storage layer. In the example below, we have `Location() tile.Point` and `SetLocation(tile.Point)` methods on our mobile entity, which contain the necessary code to read/write but the entity itself does not have any state besides a `Cursor`. This allows us to lazily evaluate and only read the necessary data when accessing it.

```go
// Mobile represents a view on a current mobile row
type Mobile struct {
	row *column.Cursor
}

// Location reads the current location
func (m *Mobile) Location() tile.Point {
    // read the "location" column
}

// SetLocation writes the current location
func (m *Mobile) SetLocation(v tile.Point) {
	// write the "location" column
}
```

## Systems

This `system` directory various game **systems** that are executed periodically and process a set of **components** (i.e. columns) for a set of **entities** (i.e. players, items, monsters). Systems access data using columnar **queries** which allow us to filter only the rows that the system can process.

For example, consider a _movement system_ that adjust both `location` and `movement action` components of a `mobile object` entity. Such system needs to first filter out entities that haven't moved, which is done using a bitmap index that checks whether the `movement action` is not empty, then the system processes all of the matching entities by updating the movement accordingly. If a lot of entities move, this is very efficient since cache misses are reduced and the movement logic is neatly contained within the system that performs the given behavior.

```go
// Update is called periodically on the movement system
func (s *System) Update(dt time.Duration) error {
	return s.mobs.Query(func(txn *column.Txn) error {
		return txn.With("moving").Range("move", func(row column.Cursor) {
			mobile := entity.MobileAt(&row)
			movement := m.Movement()
            location := m.Location()

            // Update the movement vector and store it
            movement, moved = movement.Update(dt)
            m.SetMovement(movement)
            if !moved {
                return false // not moved
            }

            // Try to move and check whether the location is within map bounds
            location = location.Move(movement.Direction())
            if !location.WithinSize(s.grid.Size){
                return false // out of map bounds
            }

            // Update the current location
            m.SetLocation(location)
		})
	})
}
```
