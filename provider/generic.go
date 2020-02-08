package provider

import (
	"math"
	"reflect"
	"sync"

	"github.com/cheekybits/genny/generic"
	"github.com/vmihailenco/msgpack/v4"
)

//go:generate genny -pkg=builtin -in=$GOFILE -out=z_components.go gen "Any=Box"

// Any is the generic type.
type Any generic.Type

// --------------------------- Component of Any ----------------------------

// TypeOfAny returns the type of the component
var TypeOfAny = reflect.TypeOf(new(Any)).Elem()

// ProviderOfAny represents an array of components.
type ProviderOfAny struct {
	sync.RWMutex
	typ  reflect.Type
	free []int
	page []pageOfAny
}

// NewProviderOfAny creates an array of components for the specific type.
func NewProviderOfAny() *ProviderOfAny {
	const cap = 128
	c := &ProviderOfAny{
		free: make([]int, 0, cap),
		page: make([]pageOfAny, 0, cap),
	}
	c.typ = reflect.TypeOf(pageOfAny{}.data).Elem()
	return c
}

// Type returns the type of the component.
func (c *ProviderOfAny) Type() reflect.Type {
	return c.typ
}

// Add adds a component to the array. Returns the index in the array which
// can be used to remove the component from the array.
func (c *ProviderOfAny) Add(v interface{}) int {
	return c.AddAny(v.(Any)) // Must be of correct type
}

// AddAny adds a component to the array. Returns the index in the array which
// can be used to remove the component from the array.
func (c *ProviderOfAny) AddAny(v Any) int {
	c.Lock()
	defer c.Unlock()

	if len(c.free) == 0 {
		pageAt := len(c.page)
		c.page = append(c.page, pageOfAny{})
		c.free = append(c.free, pageAt)
		offset := c.page[pageAt].Add(v)
		return (64 * pageAt) + offset
	}

	// find the free page and append
	last := len(c.free) - 1
	pageAt := c.free[last]
	offset := c.page[pageAt].Add(v)
	if c.page[pageAt].IsFull() {
		c.free = c.free[:last]
	}
	return (64 * pageAt) + offset
}

// View iterates over the array but only acquires a read lock. Make sure you do
// not mutate the state during this iteration as the pointer is given merely for
// performance reasons.
func (c *ProviderOfAny) View(f func(*Any)) {
	c.RLock()
	defer c.RUnlock()
	for i := 0; i < len(c.page); i++ {
		c.page[i].Range(f)
	}
}

// Update ranges over the data in the slice and lets the user update it. This
// acquires a read-write lock and is safe to update concurrently.
func (c *ProviderOfAny) Update(f func(*Any)) {
	c.Lock()
	defer c.Unlock()
	for i := 0; i < len(c.page); i++ {
		c.page[i].Range(f)
	}
}

// ViewAt returns a specific component located at the given index. Read lock
// is acquired in this operation, use it sparingly.
func (c *ProviderOfAny) ViewAt(index int) Any {
	pageAt, offset := index/64, index%64
	c.RLock()
	defer c.RUnlock()
	return *(c.page[pageAt].At(offset))
}

// UpdateAt updates a component at a specific location. Write lock is acquired
// in this operation, use it sparingly.
func (c *ProviderOfAny) UpdateAt(index int, f func(*Any)) {
	pageAt, offset := index/64, index%64
	c.Lock()
	f(c.page[pageAt].At(offset))
	c.Unlock()
}

// RemoveAt removes a component at a specific location. Write lock is acquired
// in this operation, use it sparingly.
func (c *ProviderOfAny) RemoveAt(index int) {
	pageAt, offset := index/64, index%64
	c.Lock()
	defer c.Unlock()
	if c.page[pageAt].IsFull() {
		c.free = append(c.free, pageAt)
	}
	c.page[pageAt].Del(offset)
}

// EncodeMsgpack encodes the component in message pack format into the writer.
func (c *ProviderOfAny) EncodeMsgpack(enc *msgpack.Encoder) (err error) {
	if err = enc.Encode(c.free); err == nil {
		err = enc.Encode(c.page)
	}
	return
}

// DecodeMsgpack decodes the page from the reader in message pack format.
func (c *ProviderOfAny) DecodeMsgpack(dec *msgpack.Decoder) (err error) {
	if err = dec.Decode(&c.free); err == nil {
		err = dec.Decode(&c.page)
	}
	return
}

// ---------------------------- Page of Any -----------------------------

// Page represents a page for a particular type.
type pageOfAny struct {
	full uint64
	data [64]Any
}

// Add adds an element to the page and returns the offset.
func (p *pageOfAny) Add(v Any) (index int) {
	if p.IsFull() {
		return -1
	}

	for i := 0; i < 64; i++ {
		if (p.full & (1 << i)) == 0 {
			p.full |= (1 << i)
			p.data[i] = v
			return i
		}
	}
	return -1
}

// Del deletes an element at an offset.
func (p *pageOfAny) Del(index int) {
	p.full &= uint64(^(1 << index))
}

// IsFull checks whether the page is full or not.
func (p *pageOfAny) IsFull() bool {
	return p.full == math.MaxUint64
}

// At returns a specific component located at the given index.
func (p *pageOfAny) At(index int) *Any {
	if (p.full & (1 << index)) > 0 {
		return &p.data[index]
	}
	return nil
}

// Range iterates over the page.
func (p *pageOfAny) Range(f func(*Any)) {
	if p.IsFull() {
		for i := 0; i < 64; i++ {
			f(&p.data[i])
		}
		return
	}

	for i := 0; i < 64; i++ {
		if (p.full & (1 << i)) > 0 {
			f(&p.data[i])
		}
	}
}

// Encode encodes the page in message pack format into the writer.
func (p *pageOfAny) EncodeMsgpack(enc *msgpack.Encoder) (err error) {
	if err = enc.EncodeUint64(p.full); err == nil {
		err = enc.Encode(p.data)
	}
	return
}

// Decode decodes the page from the reader in message pack format.
func (p *pageOfAny) DecodeMsgpack(dec *msgpack.Decoder) (err error) {
	if p.full, err = dec.DecodeUint64(); err == nil {
		err = dec.Decode(&p.data)
	}
	return
}
