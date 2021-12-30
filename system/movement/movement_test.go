package movement

import (
	"testing"
	"time"

	"github.com/kelindar/ecs/entity/mobile"
	"github.com/kelindar/ecs/state"
	"github.com/kelindar/ecs/world"
	"github.com/kelindar/tile"
	"github.com/stretchr/testify/assert"
)

func TestTryUpdate(t *testing.T) {
	s, _ := newSystem()

	// Move west, should be okay
	assert.NoError(t, s.mobiles.UpdateAt(0, func(v mobile.Mobile) error {
		assert.True(t, s.tryUpdate(v, time.Second))
		return nil
	}))

	// Move west again, should fail given that we reached the bounds of the map
	assert.NoError(t, s.mobiles.UpdateAt(0, func(v mobile.Mobile) error {
		assert.False(t, s.tryUpdate(v, time.Second))
		return nil
	}))
}

// newSystem creates a new system for testing purposes
func newSystem() (*System, *world.World) {
	system := new(System)
	world := world.Create(9, 9, system)
	world.Mobiles.Insert(func(v mobile.Mobile) {
		v.SetMovement(state.NewMovement(tile.West, 5, time.Second, 400*time.Millisecond))
		v.SetLocation(tile.At(1, 0))
	})
	return system, world
}
