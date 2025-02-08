package movement

import (
	"time"

	"github.com/kelindar/column"
	"github.com/kelindar/ecs/entity/mobile"
	"github.com/kelindar/ecs/state"
	"github.com/kelindar/ecs/world"
	"github.com/kelindar/tile"
)

// Assert contract compliance
var _ world.System[any] = new(System)

// System represents a system that handles all movement of mobile objects
type System struct {
	grid    *tile.Grid[any]
	mobiles *mobile.Collection
}

// Interval specifies how often the system should run
func (s *System) Interval() time.Duration {
	return 100 * time.Millisecond
}

// Attach attaches the system to the world context
func (s *System) Attach(w *world.World[any]) error {
	s.grid = w.Grid
	s.mobiles = w.Mobiles
	s.mobiles.CreateIndex("moving", "move", func(r column.Reader) bool {
		return state.Movement(r.Uint()).Distance() > 0
	})
	return nil
}

// Update is called periodically to update the system
func (s *System) Update(dt *world.Clock) error {
	elapsed := dt.Elapsed
	return s.mobiles.Range(func(m mobile.Mobile) {
		if !s.tryUpdate(m, elapsed) {
			return // No movement
		}
	}, "moving")
}

// tryUpdate attempts to update a movement state and location of the mobile
func (s *System) tryUpdate(m mobile.Mobile, dt time.Duration) (moved bool) {
	movement := m.Movement()
	location := m.Location()

	// Update the movement vector and store it
	movement, moved = movement.Update(dt)
	m.SetMovement(movement)
	if !moved {
		return false // not moved
	}

	// Try to move and check whether the location is within bounds
	location = location.Move(movement.Direction())
	_, ok := s.grid.At(location.X, location.Y)
	if !ok {
		return false // out of bounds
	}

	// TODO: check tile for terrain type

	// Update the current location
	m.SetLocation(location)
	return true
}
