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
		// Create the next world state (same dimensions as the current world)
		nextWorld := make([][]byte, p.ImageHeight)
		for i := range world {
			nextWorld[i] = make([]byte, p.ImageWidth) // Initialize each row of the new world
		}

		// Iterate over each cell in the world
		for i := 0; i < p.ImageHeight; i++ {
			for j := 0; j < p.ImageWidth; j++ {
				// Calculate the sum of the live neighbors
				sum := int(world[(i+p.ImageHeight-1)%p.ImageHeight][(j+p.ImageWidth-1)%p.ImageWidth]) +
					int(world[(i+p.ImageHeight-1)%p.ImageHeight][(j+p.ImageWidth)%p.ImageWidth]) +
					int(world[(i+p.ImageHeight-1)%p.ImageHeight][(j+p.ImageWidth+1)%p.ImageWidth]) +
					int(world[(i+p.ImageHeight)%p.ImageHeight][(j+p.ImageWidth-1)%p.ImageWidth]) +
					int(world[(i+p.ImageHeight)%p.ImageHeight][(j+p.ImageWidth+1)%p.ImageWidth]) +
					int(world[(i+p.ImageHeight+1)%p.ImageHeight][(j+p.ImageWidth-1)%p.ImageWidth]) +
					int(world[(i+p.ImageHeight+1)%p.ImageHeight][(j+p.ImageWidth)%p.ImageWidth]) +
					int(world[(i+p.ImageHeight+1)%p.ImageHeight][(j+p.ImageWidth+1)%p.ImageWidth])

				sum /= 255 // Normalize sum since alive cells are 255 (dead are 0)

				// Apply the Game of Life rules
				if world[i][j] == 255 { // Cell is alive
					if sum < 2 || sum > 3 {
						nextWorld[i][j] = 0 // Cell dies
					} else {
						nextWorld[i][j] = 255 // Cell stays alive
					}
				} else { // Cell is dead
					if sum == 3 {
						nextWorld[i][j] = 255 // Cell becomes alive
					} else {
						nextWorld[i][j] = 0 // Cell stays dead
					}
				}
			}
		}
		world = nextWorld
	}
	// TODO: report the final state using FinalTurnCompleteEvent.
	aliveCells := []util.Cell{}

	for i, row := range world {
		for j := range row {
			// If the cell is alive (i.e., has a value of 1)
			if world[i][j] == 255 {
				// Add the alive cell's coordinates (x, y) to the list of alive cells
				aliveCells = append(aliveCells, util.Cell{j, i})
			}
		}
	}
	c.events <- FinalTurnComplete{turn, aliveCells}

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
