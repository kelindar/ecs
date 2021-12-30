package snapshot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnapshot(t *testing.T) {
	var count int
	system := new(System)
	system.save = func() error {
		count++
		return nil
	}

	assert.NotZero(t, system.Interval())
	assert.NoError(t, system.Update(nil))
	assert.NoError(t, system.Close())
	assert.Equal(t, 2, count)
}
