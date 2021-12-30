package entity

import (
	"os"
	"path"

	"github.com/kelindar/column"
	"github.com/rs/xid"
)

// Collection represents a collection of mobile objects
type Collection[T any] struct {
	*column.Collection
	name string
	read func(*column.Cursor) T
}

// NewCollection creates a new mobile object collection
func NewCollection[T any](name string, read func(*column.Cursor) T) *Collection[T] {
	db := column.NewCollection()
	db.CreateColumn("id", column.ForKey()) // Unique ID
	return &Collection[T]{
		Collection: db,
		name:       name,
		read:       read,
	}
}

// Insert inserts a mobile into the collection
func (c *Collection[T]) Insert(fn func(v T)) (index uint32, err error) {
	id := xid.New().String()
	return c.Collection.Insert("id", func(v column.Cursor) error {
		v.SetString(id)
		fn(c.read(&v))
		return nil
	})
}

// Range iterates over all rows that match the specified filter columns
func (c *Collection[T]) Range(fn func(v T), filters ...string) error {
	return c.Collection.Query(func(txn *column.Txn) error {
		return txn.With(filters...).Range("id", func(v column.Cursor) {
			fn(c.read(&v))
		})
	})
}

// UpdateAt updates a mobile at a given index
func (c *Collection[T]) UpdateAt(idx uint32, fn func(v T) error) error {
	return c.Collection.UpdateAt(idx, "id", func(v column.Cursor) error {
		return fn(c.read(&v))
	})
}

// ---------------------------------- Load/Save ----------------------------------

// Restore restores the collection from the specified directory. This operation
// should be called before any of transactions, right after initialization. If
// the file does not exist, it creates an empty collection and saves it.
func (c *Collection[T]) Restore(dir string) error {
	filename := path.Join(dir, c.name)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return c.Snapshot(dir)
	}

	// Otherwise, attempt to open the file and restore
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()
	return c.Collection.Restore(file)
}

// Snapshot writes a collection snapshot into the specified directory.
func (c *Collection[T]) Snapshot(dir string) error {
	filename := path.Join(dir, c.name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	// Create a new snapshot of an empty collection
	defer file.Close()
	return c.Collection.Snapshot(file)
}
