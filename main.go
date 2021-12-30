package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelindar/ecs/system/movement"
	"github.com/kelindar/ecs/system/snapshot"
	"github.com/kelindar/ecs/world"
)

func main() {
	world, err := world.Open("save",
		new(movement.System),
		new(snapshot.System),
	)
	if err != nil {
		panic(err)
	}

	onSignal(func(_ os.Signal) {
		world.Close()
	})

	// Start the world and wait until it's stopped
	world.Simulate(context.Background())
}

// onSignal hooks a callback for a signal
func onSignal(callback func(sig os.Signal)) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range c {
			callback(sig)
		}
	}()
}
