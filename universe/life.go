package universe

import "fmt"

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
	fmt.Println("before", uni.SpaceTime[1][1])
	for i, row := range uni.SpaceTime {
		for j, _ := range row {
			uni.mu.Lock()
			uni.SpaceTime[i][j].survivesNextGen()
			uni.mu.Unlock()
		}
	}
	fmt.Println("after", uni.SpaceTime[1][1])
}

func ToNextGen() bool {
	universeDies := true
	fmt.Println("next gen")
	for i, row := range uni.SpaceTime {
		// WARNING: does this modify cells or nope?
		for j, cell := range row {
			if cell.survives {
				universeDies = false
				uni.SpaceTime[i][j].char = "*"
				uni.SpaceTime[i][j].alive = true
			} else {
				uni.SpaceTime[i][j].char = "_"
				uni.SpaceTime[i][j].alive = false
			}
		}
	}
	uni.Generation += 1
	return universeDies
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
