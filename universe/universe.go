package universe

import (
	"fmt"
	"gol/constants"
	"math/rand"
	"sync/atomic"

	"github.com/Zyko0/go-sdl3/sdl"
)

type Cell struct {
	Point    sdl.FPoint
	Alive    bool
	survives bool
}

// universe is the universe, duh
// SpaceTime it's of a fixed size, so there should be no problem
// about referencing the backing array. We'll see that :D
type Universe struct {
	SpaceTime  [constants.WINDOW_H][constants.WINDOW_W]Cell // Fixed size matrix, each element is a Cell
	Generation uint                                         // Generation number, starting from 0
	Population uint                                         // Number of alive cells, if it goes to 0, life ends
	dies       atomic.Bool
}

var uni Universe

func GetUniverse() *Universe {
	return &uni
}

// getNeighbours returns a 8 length array of all adjacent Cells, excluded the one calling the method.
func (c Cell) GetNeighbours() [8]Cell {
	return c.getNeighbours()
}
func (c Cell) getNeighbours() [8]Cell {
	n := [8]Cell{}
	index := 0
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			// Skip calling cell
			if i == 0 && j == 0 {
				continue
			}
			deltaX := int(c.Point.X) + j
			deltaY := int(c.Point.Y) + i
			// wraparound, because the Universe is infinite
			switch deltaX {
			case -1:
				deltaX = constants.WINDOW_W - 1
			case constants.WINDOW_W:
				deltaX = 0
			}
			switch deltaY {
			case -1:
				deltaY = constants.WINDOW_H - 1
			case constants.WINDOW_H:
				deltaY = 0
			}
			n[index] = uni.SpaceTime[deltaY][deltaX]
			index++
		}
	}
	return n
}

// GetCellFromUniverse returns a Cell from the Universe.
// x and y are adjusted to fit the infinite universe, meaning that if
// the value extends the grid constants, we wraparound.
func GetCellFromUniverse(x, y int) Cell {
	if x < 0 {
		tmp := -x
		x = constants.WINDOW_H - tmp
	}
	if y < 0 {
		tmp := -y
		y = constants.WINDOW_H - tmp
	}
	if x > constants.WINDOW_H-1 {
		x = x % constants.WINDOW_H
	}
	if y > constants.WINDOW_W-1 {
		y = y % constants.WINDOW_H
	}
	return uni.SpaceTime[x][y]
}

func PrintNeighbours(x, y int) {
	cell := &uni.SpaceTime[x][y]
	cell.printNeighbours()
}
func (c Cell) printNeighbours() {
	n := c.getNeighbours()
	for i := range 8 {
		fmt.Printf("[%d %d]", n[i].Point.Y, n[i].Point.X)
	}
}

func SpawnUniverse() {
	initUniverse()
	populateUniverse()
}

func initUniverse() {
	for i, row := range uni.SpaceTime {
		for j, _ := range row {
			uni.SpaceTime[i][j].Point.Y = float32(i)
			uni.SpaceTime[i][j].Point.X = float32(j)
		}
	}
}

func UpdatePopulation() {
	uni.Population = 0
	for i, row := range uni.SpaceTime {
		for j, _ := range row {
			if uni.SpaceTime[i][j].Alive {
				uni.Population++
			}
		}
	}
}

func populateUniverse() {
	alive := make([]struct{ x, y int }, constants.CELLS_ALIVE_AT_START)
	for range constants.CELLS_ALIVE_AT_START {
		randRow := rand.Intn(constants.WINDOW_H)
		randCol := rand.Intn(constants.WINDOW_W)
		alive = append(alive, struct {
			x int
			y int
		}{x: randCol, y: randRow})
	}
	for _, n := range alive {
		uni.SpaceTime[n.y][n.x].Alive = true
	}
}
