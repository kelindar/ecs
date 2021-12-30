package state

import (
	"testing"
	"time"

	"github.com/kelindar/tile"
	"github.com/stretchr/testify/assert"
)

func TestMovement(t *testing.T) {
	v := NewMovement(tile.East, 5, time.Second, 400*time.Millisecond)
	assert.Equal(t, tile.East, v.Direction())
	assert.Equal(t, 5, v.Distance())
	assert.Equal(t, time.Second, v.Velocity())
	assert.Equal(t, 400*time.Millisecond, v.Duration())
	assert.Equal(t, "movement 5🡲E, 1s/tile, 𝚫t=400ms", v.String())
}

func TestMovementUpdate(t *testing.T) {
	v := NewMovement(tile.East, 5, time.Second, 400*time.Millisecond)
	updated, moved := v.Update(time.Second)
	assert.True(t, moved)
	assert.Equal(t, 4, updated.Distance())
	assert.Equal(t, time.Second, updated.Duration())
}

func TestMovementPanic(t *testing.T) {
	assert.Panics(t, func() {
		NewMovement(tile.East, 10, time.Second, time.Second)
	})

	assert.Panics(t, func() {
		NewMovement(tile.East, 5, time.Hour, time.Hour)
	})
}
