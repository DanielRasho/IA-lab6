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

func formatCenteredCell(value string, cellWidth int) string {
	remainingSpace := cellWidth - len(value)
	padlen := (cellWidth - len(value)) / 2
	if remainingSpace%2 == 0 {
		return fmt.Sprintf("%*s%s%*s", padlen, "", value, padlen, "")
	} else {
		return fmt.Sprintf("%*s%s%*s", padlen+1, "", value, padlen, "")
	}
}

func DisplaySimultionResults(st *ReportTable) {
	fmt.Printf(
		`
====================================================================================================================
| %s | %s | %s | %s | %s |
====================================================================================================================
`,
		formatCenteredCell("Matchup", 36),
		formatCenteredCell("P1", 16),
		formatCenteredCell("P2", 16),
		formatCenteredCell("Draws", 16),
		formatCenteredCell("Game count", 16),
	)
	for players, stats := range *st {
		fmt.Printf("| %s | %s | %s | %s | %s |\n",
			formatCenteredCell(fmt.Sprintf("%s vs %s", PrettyDisplay(players.P1), PrettyDisplay(players.P2)), 36),
			formatCenteredCell(fmt.Sprintf("%f", stats.WinsP1), 16),
			formatCenteredCell(fmt.Sprintf("%f", stats.WinsP2), 16),
			formatCenteredCell(fmt.Sprintf("%f", stats.Draws), 16),
			formatCenteredCell(fmt.Sprintf("%f", stats.GameCount), 16),
		)
	}
}
