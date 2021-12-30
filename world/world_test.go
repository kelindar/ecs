package world

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorldOpen(t *testing.T) {
	defer os.RemoveAll("temp")

	{ // Create
		w, err := Open("temp")
		assert.NoError(t, err)
		assert.NotNil(t, w)
		assert.NoError(t, w.Close())
	}

	{ // Restore
		w, err := Open("temp")
		assert.NoError(t, err)
		assert.NotNil(t, w)
		assert.NoError(t, w.Close())
	}
}
