package universe

import (
	"sync"
)

/*
RULES OF GAME OF LIFE:
1) any live cell with fewer than two live neighbours dies, by underpopulation
2) any live cell with 2 or 3 live neighbours lives on to the next gen
3) any live cell with more than 3 live neighbours dies, by overpopulation
4) any dead cell with exactly 3 live neighbours becomes a live cell, by reproduction
	The first generation is created by applying the above rules simultaneously to every cell in the seed, live or dead;
	births and deaths occur simultaneously, and the discrete moment at which this happens is sometimes called a tick.
	Each generation is a pure function of the preceding one. The rules continue to be applied repeatedly to create further generations
*/

// Since writing to the universe simoultanously, has sync problems, we should create a new universe with the next gen.
// Then assign the new universe to the global uni variable.
func ApplyRules() {
	for i, row := range uni.SpaceTime {
		for j, _ := range row {
			uni.SpaceTime[i][j].survivesNextGen()
		}
	}
}

func ApplyRulesInParallel(wg *sync.WaitGroup) {
	for i, row := range uni.SpaceTime {
		for j, _ := range row {
			wg.Add(1)
			go func(ii, jj int) {
				uni.SpaceTime[ii][jj].survivesNextGen()
				wg.Done()
			}(i, j) // Closure with i, j ensures copy by value into ii, jj, allowing indipendent write to each Cell
		}
	}
}

func ToNextGen(wg *sync.WaitGroup) bool {
	uni.dies.Store(true)
	for i, row := range uni.SpaceTime {
		for j, _ := range row {
			wg.Add(1)
			go func(ii, jj int) {
				if uni.SpaceTime[ii][jj].survives {
					// dies is the only shared variable that can be affected by race condition
					// SpaceTime is also shared, true, but each Cell is indipendent, since the go routines tackle
					// one each
					uni.dies.Store(false)
					uni.SpaceTime[ii][jj].char = "*"
					uni.SpaceTime[ii][jj].alive = true
				} else {
					uni.SpaceTime[ii][jj].char = "_"
					uni.SpaceTime[ii][jj].alive = false
				}
				wg.Done()
			}(i, j)
		}
	}
	wg.Wait()
	uni.Generation += 1
	return uni.dies.Load()
}

func (c *Cell) survivesNextGen() {
	numAlive := c.getAliveNeighbours()
	if !c.alive && numAlive == 3 {
		c.survives = true
		return
	}
	if c.alive && numAlive == 2 || c.alive && numAlive == 3 {
		c.survives = true
		return
	}
	c.survives = false
}

func (c Cell) getAliveNeighbours() int {
	livingN := 0
	neighbours := c.getNeighbours()
	for _, n := range neighbours {
		if n.alive {
			livingN++
		}
	}
	return livingN
}
