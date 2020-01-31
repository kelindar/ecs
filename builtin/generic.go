package builtin

import (
	"math"
	"reflect"
	"sync"

	"github.com/cheekybits/genny/generic"
	"github.com/vmihailenco/msgpack/v4"
)

//go:generate genny -pkg=builtin -in=$GOFILE -out=z_components.go gen "TType=float32,float64,int16,int32,int64,uint16,uint32,uint64"

// TType is the generic type.
type TType generic.Type

// --------------------------- Component of TType ----------------------------

// TypeOfTType returns the type of the component
var TypeOfTType = reflect.TypeOf(new(TType)).Elem()

// ProviderOfTType represents an array of components.
type ProviderOfTType struct {
	sync.RWMutex
	typ  reflect.Type
	free []int
	page []pageOfTType
}

// NewProviderOfTType creates an array of components for the specific type.
func NewProviderOfTType() *ProviderOfTType {
	const cap = 128
	c := &ProviderOfTType{
		free: make([]int, 0, cap),
		page: make([]pageOfTType, 0, cap),
	}
	c.typ = reflect.TypeOf(pageOfTType{}.data).Elem()
	return c
}

// Type returns the type of the component.
func (c *ProviderOfTType) Type() reflect.Type {
	return c.typ
}

// Add adds a component to the array. Returns the index in the array which
// can be used to remove the component from the array.
func (c *ProviderOfTType) Add(component interface{}) int {
	v := component.(TType) // Must be of correct type
	c.Lock()
	defer c.Unlock()

	if len(c.free) == 0 {
		pageAt := len(c.page)
		c.page = append(c.page, pageOfTType{})
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
func (c *ProviderOfTType) View(f func(*TType)) {
	c.RLock()
	defer c.RUnlock()
	for i := 0; i < len(c.page); i++ {
		c.page[i].Range(f)
	}
}

// Update ranges over the data in the slice and lets the user update it. This
// acquires a read-write lock and is safe to update concurrently.
func (c *ProviderOfTType) Update(f func(*TType)) {
	c.Lock()
	defer c.Unlock()
	for i := 0; i < len(c.page); i++ {
		c.page[i].Range(f)
	}
}

// ViewAt returns a specific component located at the given index. Read lock
// is acquired in this operation, use it sparingly.
func (c *ProviderOfTType) ViewAt(index int) TType {
	pageAt, offset := index/64, index%64
	c.RLock()
	defer c.RUnlock()
	return *(c.page[pageAt].At(offset))
}

// UpdateAt updates a component at a specific location. Write lock is acquired
// in this operation, use it sparingly.
func (c *ProviderOfTType) UpdateAt(index int, f func(*TType)) {
	pageAt, offset := index/64, index%64
	c.Lock()
	f(c.page[pageAt].At(offset))
	c.Unlock()
}

// RemoveAt removes a component at a specific location. Write lock is acquired
// in this operation, use it sparingly.
func (c *ProviderOfTType) RemoveAt(index int) {
	pageAt, offset := index/64, index%64
	c.Lock()
	defer c.Unlock()
	if c.page[pageAt].IsFull() {
		c.free = append(c.free, pageAt)
	}
	c.page[pageAt].Del(offset)
}

// EncodeMsgpack encodes the component in message pack format into the writer.
func (c *ProviderOfTType) EncodeMsgpack(enc *msgpack.Encoder) (err error) {
	if err = enc.Encode(c.free); err == nil {
		err = enc.Encode(c.page)
	}
	return
}

// DecodeMsgpack decodes the page from the reader in message pack format.
func (c *ProviderOfTType) DecodeMsgpack(dec *msgpack.Decoder) (err error) {
	if err = dec.Decode(&c.free); err == nil {
		err = dec.Decode(&c.page)
	}
	return
}

// ---------------------------- Page of TType -----------------------------

// Page represents a page for a particular type.
type pageOfTType struct {
	full uint64
	data [64]TType
}

// Add adds an element to the page and returns the offset.
func (p *pageOfTType) Add(v TType) (index int) {
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
func (p *pageOfTType) Del(index int) {
	p.full &= uint64(^(1 << index))
}

// IsFull checks whether the page is full or not.
func (p *pageOfTType) IsFull() bool {
	return p.full == math.MaxUint64
}

// At returns a specific component located at the given index.
func (p *pageOfTType) At(index int) *TType {
	if (p.full & (1 << index)) > 0 {
		return &p.data[index]
	}
	return nil
}

// Range iterates over the page.
func (p *pageOfTType) Range(f func(*TType)) {
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
func (p *pageOfTType) EncodeMsgpack(enc *msgpack.Encoder) (err error) {
	if err = enc.EncodeUint64(p.full); err == nil {
		err = enc.Encode(p.data)
	}
	return
}

// Decode decodes the page from the reader in message pack format.
func (p *pageOfTType) DecodeMsgpack(dec *msgpack.Decoder) (err error) {
	if p.full, err = dec.DecodeUint64(); err == nil {
		err = dec.Decode(&p.data)
	}
	return
}
