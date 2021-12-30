package world

import (
	"reflect"
	"strings"
	"time"
)

// System represents a contract that every game system implements
type System interface {
	Interval() time.Duration
	Attach(*World) error
	Update(*Clock) error
}

// Clock represents a game clock
type Clock struct {
	Elapsed time.Duration // Elapsed time between frames
	Current time.Time     // Current time, in unix seconds
}

// newClock creates a new clock
func newClock() *Clock {
	return &Clock{
		Current: time.Now().UTC(),
	}
}

// Update updates the current clock
func (c *Clock) Update() {
	now := time.Now().UTC()
	c.Elapsed = c.Current.Sub(now)
	c.Current = now
}

// nameOf prettifies system name
func nameOf(system System) string {
	name := reflect.TypeOf(system).String()
	name = strings.TrimSuffix(name, ".System")
	name = strings.TrimPrefix(name, "*")
	return name
}
