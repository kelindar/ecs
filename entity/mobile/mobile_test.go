package mobile

import (
	"testing"
	"time"

	"github.com/kelindar/ecs/state"
	"github.com/kelindar/tile"
	"github.com/stretchr/testify/assert"
)

func TestMobile(t *testing.T) {
	c := NewCollection()
	assert.NotNil(t, c)

	// Insert
	_, err := c.Insert(func(mobile Mobile) {
		mobile.SetLocation(tile.At(1, 1))
		mobile.SetMovement(state.NewMovement(tile.East, 5, time.Second, 400*time.Millisecond))
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, c.Count())

	// Range
	assert.NoError(t, c.Range(func(mobile Mobile) {
		assert.NotEmpty(t, mobile.ID())
		assert.NotEmpty(t, mobile.Location())
		assert.NotEmpty(t, mobile.Movement())
	}))
}
