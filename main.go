package main

import (
	"fmt"
	"gol/constants"
	"gol/universe"
	"sync"
	"time"

	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/sdl"
)

/*
RULES OF GAME OF LIFE:
 1. any live cell with fewer than two live neighbours dies, by underpopulation
 2. any live cell with 2 or 3 live neighbours lives on to the next gen
 3. any live cell with more than 3 live neighbours dies, by overpopulation
 4. any dead cell with exactly 3 live neighbours becomes a live cell, by reproduction
    The first generation is created by applying the above rules simultaneously to every cell in the seed, live or dead;
    births and deaths occur simultaneously, and the discrete moment at which this happens is sometimes called a tick.
    Each generation is a pure function of the preceding one. The rules continue to be applied repeatedly to create further generations
*/
var tick = 33 * time.Millisecond

// No concurrency version

// func main() {
// 	wg := sync.WaitGroup{}
// 	universe.SpawnUniverse()
// 	universe.PrintUniverse()
//
// 	// Let's tick every 5 seconds for now
// 	// This blocks main routine until we close the channel
// 	ticker := time.NewTicker(tick)
// 	for range ticker.C {
// 		universe.ApplyRules()
// 		universe.ApplyRulesInParallel(&wg)
// 		if universe.ToNextGen() {
// 			fmt.Println("EXTINTION")
// 			break
// 		}
// 		universe.PrintUniverse()
// 	}
// }

// Concurrency version

func main1() {
	wg := sync.WaitGroup{}
	universe.SpawnUniverse()

	// Let's tick every 5 seconds for now
	// This blocks main routine until we close the channel
	ticker := time.NewTicker(tick)
	for range ticker.C {
		universe.ApplyRules()
		wg.Wait()
		if universe.ToNextGen() {
			fmt.Println("EXTINTION")
			break
		}
	}
}

const (
	WindowWidth  = 1000
	WindowHeight = 700

	NumPoints          = 500
	MinPixelsPerSecond = 30
	MaxPixelsPerSecond = 60
)

func main() {
	defer binsdl.Load().Unload() // sdl.LoadLibrary(sdl.Path())
	defer sdl.Quit()

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}

	window, renderer, err := sdl.CreateWindowAndRenderer("Game of Life", WindowWidth, WindowHeight, 0)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()
	defer window.Destroy()

	universe.SpawnUniverse()

	sdl.RunLoop(func() error {
		var event sdl.Event

		for sdl.PollEvent(&event) {
			// Close window with q
			if event.KeyboardEvent().Key == sdl.K_Q {
				return sdl.EndLoop
			}

			if event.KeyboardEvent().Key == sdl.K_N {
				universe.SpawnUniverse()
			}
			// Close window with cmd+q
			if event.Type == sdl.EVENT_QUIT {
				return sdl.EndLoop
			}
		}
		NewGame(renderer)

		return nil
	})
}

func NewGame(renderer *sdl.Renderer) {

	// Advance to next gen
	universe.ApplyRules()
	universe.ToNextGen()
	universe.UpdatePopulation()

	renderer.SetScale(1, 1)
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	renderer.SetDrawColor(255, 255, 255, 255)
	for _, row := range universe.GetUniverse().SpaceTime {
		for _, cell := range row {
			if cell.Alive {
				renderer.RenderPoint(cell.Point.X, cell.Point.Y)
			}
		}
	}

	universe.UpdatePopulation()
	renderer.SetScale(3, 3)
	renderer.SetDrawColor(0, 255, 0, 255)
	renderer.DebugText(constants.WINDOW_W-820, 0, fmt.Sprintf("Cells alive: %d", universe.GetUniverse().Population))
	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.DebugText(0, 0, fmt.Sprintf("Gen: %d", universe.GetUniverse().Generation))
	renderer.Present()
}
