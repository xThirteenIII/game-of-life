package main

import (
	"fmt"
	"gol/universe"
	"sync"
	"time"
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
func main() {
	wg := sync.WaitGroup{}
	universe.SpawnUniverse()
	universe.PrintUniverse()

	// Let's tick every 5 seconds for now
	// This blocks main routine until we close the channel
	ticker := time.NewTicker(tick)
	for range ticker.C {
		universe.ApplyRulesInParallel(&wg)
		wg.Wait()
		if universe.ToNextGen(&wg) {
			fmt.Println("EXTINTION")
			break
		}
		universe.PrintUniverse()
	}
}
