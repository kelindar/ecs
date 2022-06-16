package entity

import (
	"testing"

	"github.com/kelindar/column"
	"github.com/stretchr/testify/assert"
)

func TestCollection(t *testing.T) {
	c := NewCollection("test", cursorFor)
	c.CreateColumn("msg", column.ForString())
	assert.NotNil(t, c)

	// Insert
	err := c.Insert(func(v Object) error {
		v.SetMessage("hello")
		return nil
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
	id interface {
		Get() (string, bool)
	}
	msg interface {
		Get() (string, bool)
		Set(string)
	}
}

func cursorFor(txn *column.Txn) Object {
	return Object{
		id:  txn.Key(),
		msg: txn.String("msg"),
	}
}

func (e *Object) ID() string {
	v, _ := e.id.Get()
	return v
}

func (e *Object) Message() string {
	v, _ := e.msg.Get()
	return v
}

func (e *Object) SetMessage(v string) {
	e.msg.Set(v)
}
