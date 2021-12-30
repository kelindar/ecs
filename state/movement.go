package state

import (
	"fmt"
	"time"

	"github.com/kelindar/tile"
)

const (
	moveDelta   = 100 * time.Millisecond
	moveMaxTime = 3 * time.Second
)

// ---------------------------------- Movement ----------------------------------

// Movement represents a origin-based movement vector.
// - 3 bits representing 8 possible directions
// - 3 bits of distance, up to 7 tiles
// - 5 bits of speed (in 100ms increments, up to 3sec)
// - 5 bits of duration (in 100ms increments, up to 3sec)
type Movement uint16

// NewMovement creates a new movement vector pointing to a direction
func NewMovement(direction tile.Direction, distance int, speed, duration time.Duration) Movement {
	if distance > 7 || distance < 0 {
		panic("vector: distance must be in [0,7] range")
	}

	if duration > moveMaxTime || speed > moveMaxTime {
		panic("vector: duration and speed must be at most 3 seconds")
	}

	d := uint16(direction&0b111) << 13
	l := uint16(distance&0b111) << 10
	s := uint16(speed/moveDelta&0b11111) << 5
	t := uint16(duration / moveDelta & 0b11111)
	return Movement(d | l | s | t)
}

// Direction returns the direction (amplitude)
func (v Movement) Direction() tile.Direction {
	return tile.Direction(v >> 13)
}

// Distance returns the distance of the vector, in tiles (magnitude)
func (v Movement) Distance() int {
	return int(v >> 10 & 0b111)
}

// Velocity in milliseconds per tile, up to 3 seconds
func (v Movement) Velocity() time.Duration {
	return time.Duration(v>>5&0b11111) * moveDelta
}

// Duration returns the time left to move one tile, up to 3 seconds
func (v Movement) Duration() time.Duration {
	return time.Duration(v&0b11111) * moveDelta
}

// Update updates the movement vector based on the elapsed time and returns
// and updated vetor and whether the distance has changed or not.
func (v Movement) Update(dt time.Duration) (Movement, bool) {
	dist := v.Distance()
	left := v.Duration() - dt

	// Once the timer reaches zero, decrement a distance by one
	if left <= 0 {
		dist--
	}

	// If there's more tiles to move, update the time with the
	// velocity
	if dist > 0 {
		left = v.Velocity()
	}

	return NewMovement(v.Direction(), dist, v.Velocity(), left), dist != v.Distance()
}

// String returns string representation of a movement vector, for debugging
func (v Movement) String() string {
	return fmt.Sprintf("movement %d%s, %s/tile, 𝚫t=%s", v.Distance(), v.Direction(), v.Velocity(), v.Duration())
}
