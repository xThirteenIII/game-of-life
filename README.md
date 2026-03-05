# game-of-life

Reading from [Wikipedia](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life), I'm implementing a version of the Game of Life in Go.  
No Claude or LLMs used for coding. Used it for benchmarks and clarify some concepts.
Just fun!

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
The universe is a \[ROW_NUM\]\[COL_NUM]Cell matrix. 25x100 to fix the terminal size, 100x100 for benchmark memory and speed checks.  
Should I make a copy of the universe each generation, or modify the existing one directly?  
For now, I'm opting for a `survived` variable in the `Cell` struct, that tells if the current Cell has to survive to the next gen, after applying all the rules.
For the tick which dictates the discrete time where births and deaths of new gen occur, I use `time.NewTicker()`, which conviniently serves the purpose.  
Adding a mutex to the Universe, to ensure thread-safety writing, even though I'm not sure it's needed yet.
Each tick, i should apply the rules to every cell, each in his separate goroutine, to simulate simultaneous behaviour.

## Sequential functions

The program works fine with no concurrency. In the `main` package, we listen on the ticker.C channel. Every 33 milliseconds, a tick is sent on the channel (that is ~30 fps).  
Each tick, following things happen in order:
1. Rules are applied. For every cell:
    - call `survivesNextGen()`
    - call `getNeighbours()` to get an 8-size array with each neighbour Cell
    - set `survives: true/false` based on rules
2. Compute `ToNextGen`:
    - set `universeDies: true`, a variable which tells if extintion is reached or not (no more cells can reproduce)
    - for each cell, if it `survives: true` make it alive and set `universeDies: false`
    - or kill the cell if `survives: false`
2. Print new universe

## Concurrent go routines

In theory, the use of go routines should speed up time per op, that is if the universe size is big enough. Otherwise I found out the overhead of the  
go routines is just not worth it and the sequential program wins over go concurrency.
That said, I made `ApplyRules` to be concurrent, so each go routine can work on a single cell at the same time.
Same thing goes for `ToNextGen()`.
I'm using the `benchmark` Go package to test both sequential and concurrent programs.

#### Test: `go test ./universe/ -bench=. -benchtime=10s` Universe: 100x100
```
goos: linux
goarch: amd64
pkg: gol/universe
cpu: AMD Ryzen AI 9 365 w/ Radeon 880M              
BenchmarkApplyRulesSequential-20    	   36598	    325622 ns/op  ~0.3 ms
BenchmarkApplyRulesConcurrent-20    	    6062	   2455920 ns/op  ~2.5 ms
PASS
ok  	gol/universe	26.815s
```
So basically 100x100 is too little to see concurrency win over sequential. The header, runtime scheduler, wg sync etc. is too much.  
Also, each iteration doesn't have a lot of work to do, just `getNeighbours()`, which is really fast anyways.
Even oging up to 1000x1000 produces the same effect.  
#### Solution
Worker pools. That is, I use 1 go routine on top of each core of the 20 my AMD Ryzen has. This way, each goroutine can process multiple cells in parallel.

#### CPU Bound
So I mistakenly thought that concurrency = speed. That's not the case. The Game of Life has CPU operations, not I/O like web-servers, network APIs, reading from disk, where your process has actually dead time (idle) while waiting for data.
So actually using go routines is a problem, since I'm basically forcing my processes to do context switch, scheduling, on operations that are very fast since it's the CPU that does them :/ 

#### mutex vs atomic
If like in my case I need to lock just one variable in my data struct (universeDies), I can avoid creating a mutex (heavy) that locks the whole struct,  
and just use an atomic.Bool variable, that does the job better since by design it's just a 1 CPU instruction, and thus can't be partially accessed by other routines.

