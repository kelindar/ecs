package item

import (
	"testing"

	"github.com/kelindar/tile"
	"github.com/stretchr/testify/assert"
)

func TestItem(t *testing.T) {
	c := NewCollection()
	assert.NotNil(t, c)

	// Insert
	_, err := c.Insert(func(v Item) {
		v.SetLocation(tile.At(1, 1))
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, c.Count())

	// Range
	assert.NoError(t, c.Range(func(v Item) {
		assert.NotEmpty(t, v.ID())
		assert.NotEmpty(t, v.Location())
	}))
}
