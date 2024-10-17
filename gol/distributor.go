package gol

import "uk.ac.bris.cs/gameoflife/util"

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	// create a 2D slice to store the world.
	world := make([][]uint8, p.ImageHeight)
	for y := 0; y < p.ImageHeight; y++ {
		world[y] = make([]uint8, p.ImageWidth)
	}

	turn := 0
	c.events <- StateChange{turn, Executing}

	// TODO: complete all turns for GOL

	for i := 0; i < p.ImageHeight; i++ {
		for j := 0; j < p.ImageWidth; j++ {
			world[i][j] = <-c.ioInput //creating the initial world
		}
	}
	
	for turn < p.Turns {

	}
	// TODO: report the final state using FinalTurnCompleteEvent.
	aliveCells := []util.Cell{}
	c.events <- FinalTurnComplete{turn, aliveCells}

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
