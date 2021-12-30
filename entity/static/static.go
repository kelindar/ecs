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
	db := entity.NewCollection("statics.bin", At)
	db.CreateColumn("img", column.ForUint32()) // Image index
	db.CreateColumn("at", column.ForUint32())  // Location as packed tile.Point
	return db
}

// Static represents a view on a current row
type Static struct {
	row *column.Cursor
}

// At reads the static object entity at the cursor
func At(row *column.Cursor) Static {
	return Static{row: row}
}

// ID returns the unique identifier of the static object
func (e *Static) ID() string {
	return e.row.StringAt("id")
}

// ---------------------------------- Location ----------------------------------

// Location reads the current location
func (e *Static) Location() tile.Point {
	at := e.row.UintAt("at")
	return tile.At(int16(at>>16), int16(at))
}

// SetLocation writes the current location
func (e *Static) SetLocation(v tile.Point) {
	e.row.SetUint32At("at", v.Integer())
}
