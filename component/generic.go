package component

import (
	"math"
	"reflect"
	"sync"

	"github.com/cheekybits/genny/generic"
	"github.com/kelindar/ecs"
	"github.com/vmihailenco/msgpack"
)

//go:generate genny -in=$GOFILE -out=z_components.go gen "TType=string,bool,float32,float64,int16,int32,int64,uint16,uint32,uint64,Vector2,Vector3"

// TType is the generic type.
type TType generic.Type

// --------------------------- Component of TType ----------------------------

// OfTType represents an array of components.
type OfTType struct {
	sync.RWMutex
	typ  relfect.Type
	free []int
	page []pageOfTType
}

// ForTType creates an array of components for the specific type.
func ForTType() *OfTType {
	const cap = 128
	c := &OfTType{
		free: make([]int, 0, cap),
		page: make([]pageOfTType, 0, cap),
	}
	c.typ = relfect.TypeOf(c)
	return c
}

// Type returns the type of the component.
func (c *OfTType) Type() reflect.Type {
	return c.typ
}

// Add adds a component to the array. Returns the index in the array which
// can be used to remove the component from the array.
func (c *OfTType) Add(entity *ecs.Entity, v TType) {
	c.Lock()
	defer c.Unlock()

	if len(c.free) == 0 {
		pageAt := len(c.page)
		c.page = append(c.page, pageOfTType{})
		c.free = append(c.free, pageAt)
		c.attach(entity, pageAt, c.page[pageAt].Add(v))
		return
	}

	// find the free page and append
	last := len(c.free) - 1
	pageAt := c.free[last]
	offset := c.page[pageAt].Add(v)
	if c.page[pageAt].IsFull() {
		c.free = c.free[:last]
	}
	c.attach(entity, pageAt, offset)
}

// attach attaches the remove function to the entity.
func (c *OfTType) attach(entity *ecs.Entity, pageAt, offset int) {
	index := (64 * pageAt) + offset
	entity.Attach(func() {
		c.Lock()
		defer c.Unlock()
		pageAt, offset := index/64, index%64
		if c.page[pageAt].IsFull() {
			c.free = append(c.free, pageAt)
		}
		c.page[pageAt].Del(offset)
	})
}

// View iterates over the array but only acquires a read lock. Make sure you do
// not mutate the state during this iteration as the pointer is given merely for
// performance reasons.
func (c *OfTType) View(f func(*TType)) {
	c.RLock()
	defer c.RUnlock()
	for i := 0; i < len(c.page); i++ {
		c.page[i].Range(f)
	}
}

// Update ranges over the data in the slice and lets the user update it. This
// acquires a read-write lock and is safe to update concurrently.
func (c *OfTType) Update(f func(*TType)) {
	c.Lock()
	defer c.Unlock()
	for i := 0; i < len(c.page); i++ {
		c.page[i].Range(f)
	}
}

// EncodeMsgpack encodes the component in message pack format into the writer.
func (c *OfTType) EncodeMsgpack(enc *msgpack.Encoder) (err error) {
	if err = enc.Encode(c.free); err == nil {
		err = enc.Encode(c.page)
	}
	return
}

// DecodeMsgpack decodes the page from the reader in message pack format.
func (c *OfTType) DecodeMsgpack(dec *msgpack.Decoder) (err error) {
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
