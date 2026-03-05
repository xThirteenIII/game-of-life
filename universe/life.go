package universe

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
func (u *Universe) ToNextGen() {

}
