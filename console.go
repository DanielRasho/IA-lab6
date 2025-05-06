package main

import (
	"fmt"

	sim "github.com/DanielRasho/IA-lab6/simulation"
)

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func DisplayBoard(g sim.Game) {
	for i := range 3 {
		fmt.Printf(
			"| %s | %s | %s |\n",
			sim.CellToString(g.Board[i*3], int64(i*3+1)),
			sim.CellToString(g.Board[i*3+1], int64(i*3+2)),
			sim.CellToString(g.Board[i*3+2], int64(i*3+3)),
		)
	}
}

func DisplayReport(st sim.GameStats) {
	fmt.Printf(
		`
==========
GAME STATS
==========
Winner     : %s
# Movements: %d
Duration   : %s
`, st.Winner, st.MovementCount, st.EndTime.Sub(st.StartTime))
}
