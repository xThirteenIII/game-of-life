package universe

import (
	"fmt"
	"gol/constants"
	"math/rand"
	"sync"
)

type Cell struct {
	Pos      pos
	char     string
	alive    bool
	survives bool
}

type pos struct {
	row int // row is y
	col int // col is x
}

// universe is the universe, duh
// SpaceTime it's of a fixed size, so there should be no problem
// about referencing the backing array. We'll see that :D
type Universe struct {
	SpaceTime  [constants.ROW_NUM][constants.COL_NUM]Cell // Fixed size matrix, each element is a Cell
	Generation uint                                       // Generation number, starting from 0
	Population uint                                       // Number of alive cells, if it goes to 0, life ends
	mu         sync.RWMutex                               // mutex for thread-safety write to SpaceTime
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
			deltaX := c.Pos.col + j
			deltaY := c.Pos.row + i
			// wraparound, because the Universe is infinite
			switch deltaX {
			case -1:
				deltaX = constants.COL_NUM - 1
			case constants.COL_NUM:
				deltaX = 0
			}
			switch deltaY {
			case -1:
				deltaY = constants.ROW_NUM - 1
			case constants.ROW_NUM:
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
		x = constants.ROW_NUM - tmp
	}
	if y < 0 {
		tmp := -y
		y = constants.ROW_NUM - tmp
	}
	if x > constants.ROW_NUM-1 {
		x = x % constants.ROW_NUM
	}
	if y > constants.COL_NUM-1 {
		y = y % constants.ROW_NUM
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
		fmt.Printf("[%d %d]", n[i].Pos.row, n[i].Pos.col)
	}
}

func SpawnUniverse() {
	initUniverse()
	populateUniverse()
	uni.SpaceTime[1][1].alive = true
	uni.SpaceTime[1][1].char = "*"
	uni.SpaceTime[1][1].survives = true
}

func initUniverse() {
	for i, row := range uni.SpaceTime {
		for j, _ := range row {
			uni.SpaceTime[i][j].char = "_"
			uni.SpaceTime[i][j].Pos = pos{row: i, col: j}
		}
	}
}

func populateUniverse() {
	alive := make([]pos, constants.CELLS_ALIVE_AT_START)
	for range constants.CELLS_ALIVE_AT_START {
		randRow := rand.Intn(constants.ROW_NUM)
		randCol := rand.Intn(constants.COL_NUM)
		alive = append(alive, pos{col: randCol, row: randRow})
	}
	for _, n := range alive {
		uni.SpaceTime[n.row][n.col].char = "*"
		uni.SpaceTime[n.row][n.col].alive = true
	}
}

func PrintUniverse() {
	printUniverse()
}
func printUniverse() {
	fmt.Println("---------------------------GEN", uni.Generation, "---------------------------")
	fmt.Printf("\n\n")

	for i, row := range uni.SpaceTime {
		for j, _ := range row {
			fmt.Printf("%s", uni.SpaceTime[i][j].char)
		}
		fmt.Printf("\n")
	}
}
