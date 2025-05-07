package main

import (
	"fmt"

	sim "github.com/DanielRasho/IA-lab6/simulation"
)

type Color = int

const (
	BLACK Color = iota
	RED
	GREEN
	YELLOW
	BLUE
	MAGENTA
	CYAN
	WHITE
)

type Style = int

const (
	Standard_Text       Style = 3
	Standard_Background Style = 4
	Bright_Text         Style = 9
	Bright_Background   Style = 10
)

const RESET = "\x1b[0m"

func S(content string, color Color, bgfg Style) string {
	return fmt.Sprintf("\x1b[%d%dm%s%s", bgfg, color, content, RESET)
}

func StyleIfPlayer(mark sim.CellMark, i int64) string {
	switch mark {
	case sim.X:
		return S(sim.CellToString(mark, i), GREEN, Standard_Text)
	case sim.O:
		return S(sim.CellToString(mark, i), BLUE, Standard_Text)
	default:
		return sim.CellToString(mark, i)
	}
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func DisplayBoard(g sim.Game) {
	for i := range 3 {
		fmt.Printf(
			"| %s | %s | %s |\n",
			StyleIfPlayer(g.Board[i*3], int64(i*3+1)),
			StyleIfPlayer(g.Board[i*3+1], int64(i*3+2)),
			StyleIfPlayer(g.Board[i*3+2], int64(i*3+3)),
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
