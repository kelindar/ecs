package snapshot

import (
	"time"

	"github.com/kelindar/ecs/world"
)

// Assert contract compliance
var _ world.System = new(System)

// System represents a system that handles all movement of mobile objects
type System struct {
	save func() error
}

// Interval specifies how often the system should run
func (s *System) Interval() time.Duration {
	return 60 * time.Second
}

// Attach attaches the system to the world context
func (s *System) Attach(w *world.World) error {
	s.save = w.Save
	return nil
}

// Update is called periodically to update the system
func (s *System) Update(dt *world.Clock) error {
	return s.save()
}

// Close closes the system gracefully
func (s *System) Close() error {
	return s.save()
}
