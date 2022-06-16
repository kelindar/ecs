package static

import (
	"github.com/kelindar/column"
	"github.com/kelindar/ecs/entity"
	"github.com/kelindar/tile"
)

// Collection represents a collection of static objects
type Collection = entity.Collection[Static]

// NewCollection creates a new mobile object collection
func NewCollection() *Collection {
	db := entity.NewCollection("statics.bin", fromTxn)
	db.CreateColumn("img", column.ForUint32()) // Image index
	db.CreateColumn("at", column.ForUint32())  // Location as packed tile.Point
	return db
}

// Static represents a view on a current row
type Static struct {
	id interface {
		Get() (string, bool)
	}
	img interface {
		Get() (uint32, bool)
		Set(value uint32)
	}
	at interface {
		Get() (uint32, bool)
		Set(value uint32)
	}
}

// fromTxn creates a statically-typed mapping for a transaction
func fromTxn(txn *column.Txn) Static {
	return Static{
		id:  txn.Key(),
		at:  txn.Uint32("at"),
		img: txn.Uint32("img"),
	}
}

// ID returns the unique identifier of the item
func (e *Static) ID() string {
	v, _ := e.id.Get()
	return v
}

// ---------------------------------- Location ----------------------------------

// Location reads the current location
func (e *Static) Location() tile.Point {
	at, _ := e.at.Get()
	return tile.At(int16(at>>16), int16(at))
}

// SetLocation writes the current location
func (e *Static) SetLocation(v tile.Point) {
	e.at.Set(v.Integer())
}
