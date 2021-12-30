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
	db := entity.NewCollection("mobiles.bin", At)
	db.CreateColumn("img", column.ForUint32())  // Image index
	db.CreateColumn("at", column.ForUint32())   // Location as packed tile.Point
	db.CreateColumn("move", column.ForUint16()) // Movement vector
	return db
}

// Mobile represents a view on a current row
type Mobile struct {
	row *column.Cursor
}

// At reads the mobile object entity at the cursor
func At(row *column.Cursor) Mobile {
	return Mobile{row: row}
}

// ID returns the unique identifier of the mobile object
func (e *Mobile) ID() string {
	return e.row.StringAt("id")
}

// ---------------------------------- Location ----------------------------------

// Location reads the current location
func (e *Mobile) Location() tile.Point {
	at := e.row.UintAt("at")
	return tile.At(int16(at>>16), int16(at))
}

// SetLocation writes the current location
func (e *Mobile) SetLocation(v tile.Point) {
	e.row.SetUint32At("at", v.Integer())
}

// ---------------------------------- Movement ----------------------------------

// Movement reads the movement action
func (e *Mobile) Movement() state.Movement {
	return state.Movement(e.row.UintAt("move"))
}

// SetMovement writes the movement action
func (e *Mobile) SetMovement(v state.Movement) {
	e.row.SetUint16At("move", uint16(v))
}
