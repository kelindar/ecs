package entity

import (
	"testing"

	"github.com/kelindar/column"
	"github.com/stretchr/testify/assert"
)

func TestCollection(t *testing.T) {
	c := NewCollection("test", func(c *column.Cursor) Object {
		return Object{c}
	})
	c.CreateColumn("msg", column.ForString())
	assert.NotNil(t, c)

	// Insert
	_, err := c.Insert(func(v Object) {
		v.SetMessage("hello")
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, c.Count())

	// Update
	assert.NoError(t, c.UpdateAt(0, func(v Object) error {
		v.SetMessage("hi")
		return nil
	}))

	// Range
	assert.NoError(t, c.Range(func(v Object) {
		assert.NotEmpty(t, v.ID())
		assert.Equal(t, v.Message(), "hi")
	}))

}

// ---------------------------------- Test object ----------------------------------

type Object struct {
	row *column.Cursor
}

func (e *Object) ID() string {
	return e.row.StringAt("id")
}

func (e *Object) Message() string {
	return e.row.StringAt("msg")
}

func (e *Object) SetMessage(v string) {
	e.row.SetStringAt("msg", v)
}
