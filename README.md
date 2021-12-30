# Example of Entity Component System in Go

## Introduction

This is my attempt to build a high-performance ECS (Entity, Component, System) in pure Go. Having investigated other ECS systems which have been written to date, I realised that none of them is really a pure ECS as the data layout is often using pointers or interfaces, which would discard most of the benefits of such a system.

Instead of being an ECS framework, this repository contains an example of using [kelindar/column](https://github.com/kelindar/column) as the underlying storage of components and [kelindar/tile](https://github.com/kelindar/tile) as a 2D grid engine for spatial queries.

## Components

The components are created using a columnar storage [kelindar/column](https://github.com/kelindar/column). The library organizes data in dense arrays with bitmaps and indexing for querying, allowing us to build systems on top which simply perform queries to and update different columns when necessary.

In order to simplify the logic a bit, `entity.Collection[T]` structure wraps around `column.Collection` to provide strong typing for working with our **entities**.

```go
// Collection represents a collection of mobile objects
type Collection = entity.Collection[Mobile]

// NewCollection creates a new mobile object collection
func NewCollection() *Collection {
	db := entity.NewCollection("mobiles.bin", At)
	db.CreateColumn("img", column.ForUint32())  // Image index
	db.CreateColumn("at", column.ForUint32())   // Location as packed tile.Point
	db.CreateColumn("move", column.ForUint16()) // Movement vector
	return db
}
```

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
// Interval specifies how often the system should run
func (s *System) Interval() time.Duration {
	return 100 * time.Millisecond
}

// Attach attaches the system to the world context
func (s *System) Attach(w *world.World) error {
    s.grid = w.Grid

    // Create an index "moving" which will filter only entities
    // that have "move" field with a distance greater than zero.
    s.mobiles = w.Mobiles
    s.mobiles.CreateIndex("moving", "move", func(r column.Reader) bool {
        return state.Movement(r.Uint()).Distance() > 0
    })
    return nil
}

// Update is called periodically on the movement system
func (s *System) Update(dt time.Duration) error {
	return s.mobiles.Range(func(m mobile.Mobile) {
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
		}, "moving") // use moving index
	})
}
```
