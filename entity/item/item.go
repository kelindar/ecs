package item

import (
	"github.com/kelindar/column"
	"github.com/kelindar/ecs/entity"
	"github.com/kelindar/tile"
)

// Collection represents a collection of items
type Collection = entity.Collection[Item]

// NewCollection creates a new mobile object collection
func NewCollection() *Collection {
	db := entity.NewCollection("items.bin", fromTxn)
	db.CreateColumn("img", column.ForUint32()) // Image index
	db.CreateColumn("at", column.ForUint32())  // Location as packed tile.Point
	return db
}

// fromTxn creates a statically-typed mapping for a transaction
func fromTxn(txn *column.Txn) Item {
	return Item{
		id: txn.Key(),
		at: txn.Uint32("at"),
	}
}

// Item represents a view on a current row
type Item struct {
	id interface {
		Get() (string, bool)
	}
	at interface {
		Get() (uint32, bool)
		Set(value uint32)
	}
}

// ID returns the unique identifier of the item
func (e *Item) ID() string {
	v, _ := e.id.Get()
	return v
}

// ---------------------------------- Location ----------------------------------

// Location reads the current location
func (e *Item) Location() tile.Point {
	at, _ := e.at.Get()
	return tile.At(int16(at>>16), int16(at))
}

// SetLocation writes the current location
func (e *Item) SetLocation(v tile.Point) {
	e.at.Set(v.Integer())
}
