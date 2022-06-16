package world

import (
	"context"
	"io"
	"log"
	"runtime/debug"
	"sync"
	"time"

	"github.com/kelindar/ecs/entity/item"
	"github.com/kelindar/ecs/entity/mobile"
	"github.com/kelindar/ecs/entity/static"
	"github.com/kelindar/tile"
	"go.uber.org/multierr"
)

// World represents the entire game world state
type World struct {
	path    string             // The directory for save files
	cancel  context.CancelFunc // Cancel function to stop everything
	threads sync.WaitGroup     // Signals for each running system
	systems []System           // Attached systems
	Grid    *tile.Grid         // 3072x3072 map
	Mobiles *mobile.Collection // List of mobiles (NPCs, Players, Monsters, ...)
	Statics *static.Collection // List of objects on the map (Buildings, Trees, ...)
	Items   *item.Collection   // List of items off map (Weapons, Potions, ...)
}

// Open opens the world state file, or creates a new one
func Open(path string, systems ...System) (*World, error) {
	world := Create(3072, 3072) // Do not specify systems
	world.path = path

	// Load or create all of the collections
	if err := multierr.Combine(
		world.Mobiles.Restore(path),
		world.Statics.Restore(path),
		world.Items.Restore(path),
	); err != nil {
		return nil, err
	}

	// Register all of the provided systems
	if err := world.register(systems); err != nil {
		return nil, err
	}

	return world, nil
}

// Create creates a new empty world
func Create(width, height int16, systems ...System) *World {
	world := &World{
		Grid:    tile.NewGrid(width, height),
		Mobiles: mobile.NewCollection(),
		Statics: static.NewCollection(),
		Items:   item.NewCollection(),
	}

	// If systems are specified, attach them right away
	if err := world.register(systems); err != nil {
		panic(err)
	}
	return world
}

// Save saves the state of the world
func (w *World) Save() error {
	defer func(start time.Time) {
		log.Printf("world: save completed (%v)", time.Now().Sub(start))
	}(time.Now())
	return multierr.Combine(
		w.Mobiles.Snapshot(w.path),
		w.Statics.Snapshot(w.path),
		w.Items.Snapshot(w.path),
	)
}

// Register registers all of the systems and atta
func (w *World) register(systems []System) error {
	w.systems = systems
	for _, system := range systems {
		log.Printf("world: attaching %v system", nameOf(system))
		if err := system.Attach(w); err != nil {
			return err
		}
	}
	return nil
}

// Simulate runs the world simulation loop by starting all of the registered
// systems asynchronously.
func (w *World) Simulate(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	w.cancel = cancel

	// Every system will run in a separate goroutine
	w.threads.Add(len(w.systems))
	for _, system := range w.systems {
		go w.runSystem(ctx, system)
	}

	log.Printf("world: started successfully")
	w.threads.Wait()
	return nil
}

// runSystem runs a system
func (w *World) runSystem(ctx context.Context, system System) {
	clock := newClock()
	timer := time.NewTicker(system.Interval())

	// Update wraps the system update method and a panic handler
	update := func() {
		defer handlePanic()
		if err := system.Update(clock); err != nil {
			log.Printf("error: %+v", err)
		}
	}

	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			w.threads.Done()
			return
		case <-timer.C:
			update()
		}
	}
}

// Close saves the state of the world and closes it
func (w *World) Close() error {
	w.threads.Add(1) // Wait for closing
	if w.cancel != nil {
		w.cancel()
	}

	// Wait for all systems to stop updating, then attempt to close them
	for _, system := range w.systems {
		if closer, ok := system.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				log.Printf("world: unable to close %T, %+v", system, err)
			}
		}
	}

	// Done closing
	w.threads.Done()
	w.threads.Wait()
	return nil
}

// handlePanic handles the panic and logs it out
func handlePanic() {
	if r := recover(); r != nil {
		log.Printf("panic: %s \n %s", r, debug.Stack())
	}
}
