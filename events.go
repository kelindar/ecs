// Copyright (c) Roman Atachiants and contributore. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for detaile.

package ecs

import (
	"context"
	"sync"
)

// Events represents an event bus used by the manager.
type events struct {
	lock sync.RWMutex
	subs map[string][]*handler
}

// Notify notifies listeners of an event that happened
func (e *events) Notify(event string, value interface{}) {
	e.lock.RLock()
	defer e.lock.RUnlock()

	if handlers, ok := e.subs[event]; ok {
		for _, h := range handlers {
			h.buffer <- value
		}
	}
}

// On registers an event listener on a system
func (e *events) On(event string, callback func(interface{})) context.CancelFunc {
	e.lock.Lock()
	defer e.lock.Unlock()

	// If we don't have any listener, setup
	if e.subs == nil {
		e.subs = make(map[string][]*handler, 8)
	}

	// Create the handler
	ctx, cancel := context.WithCancel(context.Background())
	subscriber := &handler{
		buffer:   make(chan interface{}, 1),
		callback: &callback,
		cancel:   cancel,
	}

	// Add the listener
	e.subs[event] = append(e.subs[event], subscriber)
	go subscriber.listen(ctx)

	return e.unsubscribe(event, &callback)
}

// unsubscribe deregisters an event listener from a system
func (e *events) unsubscribe(event string, callback *func(interface{})) context.CancelFunc {
	return func() {
		e.lock.Lock()
		defer e.lock.Unlock()

		if handlers, ok := e.subs[event]; ok {
			clean := make([]*handler, 0, len(handlers))
			for _, h := range handlers {
				if h.callback != callback { // Compare address
					clean = append(clean, h)
				} else {
					h.cancel()
				}
			}
		}
	}
}

// -------------------------------------------------------------------------------------------

type handler struct {
	buffer   chan interface{}
	callback *func(interface{})
	cancel   context.CancelFunc
}

// Listen listens on the buffer and invokes the callback
func (h *handler) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case value := <-h.buffer:
			(*h.callback)(value)
		}
	}
}
