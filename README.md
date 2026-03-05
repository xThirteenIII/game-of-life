# game-of-life

Reading from [Wikipedia](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life), I'm implementing a version of the Game of Life in Go.  
No Claude or LLMs used. Just fun stuff!  

## Rules
The universe of the Game of Life is an infinite, two-dimensional orthogonal grid of square cells, each of which is in one of two possible states,  
live or dead (or populated and unpopulated, respectively). Every cell interacts with its eight neighbours, which are the cells that are horizontally, vertically, or diagonally adjacent.  
At each step in time, the following transitions occur:  
- Any live cell with fewer than two live neighbours dies, as if by underpopulation.
- Any live cell with two or three live neighbours lives on to the next generation.
- Any live cell with more than three live neighbours dies, as if by overpopulation.
- Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.  
The initial pattern constitutes the seed of the system. The first generation is created by applying the above rules simultaneously to every cell in the seed, live or dead; births and deaths occur simultaneously, and the discrete moment at which this happens is sometimes called a tick

## My take while developing
I'm trying to implement stuff one at a time and verifying everything works correctly. As for now, the decisions to be made are about whether to have 
methods or pure functions, whether to use value or reference methods and mainly, how to implement the simoultanious calls to every rule.  
Should I make a copy of the universe each generation, or modify the existing one directly?  
For now, I'm opting for a `survived` variable in the `Cell` struct, that tells if the current Cell has to survive to the next gen, after applying all the rules.
For the tick which dictates the discrete time where births and deaths of new gen occur, I use `time.NewTicker()`, which conviniently serves the purpose.  
Adding a mutex to the Universe, to ensure thread-safety writing, even though I'm not sure it's needed yet.
Each tick, i should apply the rules to every cell, each in his separate goroutine, to simulate simultaneous behaviour.
