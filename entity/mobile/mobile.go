package mobile

import (
	"github.com/kelindar/column"
	"github.com/kelindar/ecs/entity"
	"github.com/kelindar/ecs/state"
	"github.com/kelindar/tile"
)

// Collection represents a collection of mobile objects
type Collection = entity.Collection[Mobile]

// NewCollection creates a new mobile object collection
func NewCollection() *Collection {
	db := entity.NewCollection("mobiles.bin", fromTxn)
	db.CreateColumn("img", column.ForUint32())  // Image index
	db.CreateColumn("at", column.ForUint32())   // Location as packed tile.Point
	db.CreateColumn("move", column.ForUint16()) // Movement vector
	return db
}

// Mobile represents a view on a current row
type Mobile struct {
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
	move interface {
		Get() (uint16, bool)
		Set(value uint16)
	}
}

// fromTxn creates a statically-typed mapping for a transaction
func fromTxn(txn *column.Txn) Mobile {
	return Mobile{
		id:   txn.Key(),
		at:   txn.Uint32("at"),
		img:  txn.Uint32("img"),
		move: txn.Uint16("move"),
	}
}

// ID returns the unique identifier of the item
func (e *Mobile) ID() string {
	v, _ := e.id.Get()
	return v
}

// ---------------------------------- Location ----------------------------------

// Location reads the current location
func (e *Mobile) Location() tile.Point {
	at, _ := e.at.Get()
	return tile.At(int16(at>>16), int16(at))
}

// SetLocation writes the current location
func (e *Mobile) SetLocation(v tile.Point) {
	e.at.Set(v.Integer())
}

// ---------------------------------- Movement ----------------------------------

// Movement reads the movement action
func (e *Mobile) Movement() state.Movement {
	move, _ := e.move.Get()
	return state.Movement(move)
}

// SetMovement writes the movement action
func (e *Mobile) SetMovement(v state.Movement) {
	e.move.Set(uint16(v))
}
