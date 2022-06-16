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
	err := c.Insert(func(v Item) error {
		v.SetLocation(tile.At(1, 1))
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, c.Count())

	// Range
	assert.NoError(t, c.Range(func(v Item) {
		assert.NotEmpty(t, v.ID())
		assert.NotEmpty(t, v.Location())
	}))
}
