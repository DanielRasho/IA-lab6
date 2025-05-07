package simulation

import (
	"time"
)

type TicTacToePlayer interface {
	// Marks a specific cell on the specified board!
	// YOU MUST RETURN A NUMBER BETWEEN 0-8, this tells the system which cell you're trying to mark!
	// If the cell can't be marked because it already is, the whole world explodes!
	MakeMove(board TicTacToeBoard, whoami Turn) int
}

type GameStats struct {
	StartTime     time.Time
	EndTime       time.Time
	MovementCount uint
	Winner        string
}

type Turn = int

const (
	P1 Turn = 1
	P2 Turn = 2
)

func GetOpponent(whoami Turn) Turn {
	if whoami == P1 {
		return P2
	} else {
		return P1
	}
}

func GetMarks(whoami Turn) (CellMark, CellMark) {
	if whoami == P1 {
		return X, O
	} else {
		return O, X
	}
}

type Game struct {
	Player1 TicTacToePlayer
	Player2 TicTacToePlayer

	CurrentTurn Turn
	Board       TicTacToeBoard

	filledCells uint
	started     bool
	shouldEnd   bool
	stats       GameStats
}

func NewGameWith(p1 TicTacToePlayer, p2 TicTacToePlayer) Game {
	return Game{
		Player1:     p1,
		Player2:     p2,
		CurrentTurn: P1,
		Board:       make([]CellMark, 9),
	}
}

func (g *Game) Start() {
	g.stats.StartTime = time.Now()
	g.started = true
}

func (g Game) ShouldEnd() bool {
	return g.shouldEnd
}

func PlayerWonInBoard(b *TicTacToeBoard, player Turn) bool {
	patterns := [][]int{
		// Horizontal lines
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},

		// Vertical lines
		{0, 3, 6},
		{1, 4, 7},
		{2, 5, 8},

		// Diagonal lines
		{0, 4, 8},
		{6, 4, 2},
	}
	board := *b

patternLoop:
	for _, pattern := range patterns {
		firstCellMark := board[pattern[0]]
		if firstCellMark == EMPTY {
			continue
		}

		for i := 1; i < 3; i++ {
			if board[pattern[i]] != firstCellMark {
				continue patternLoop
			}
		}

		if firstCellMark == X && player == P1 {
			return true
		} else {
			return false
		}
	}

	return false
}

func (g *Game) Tick() {
	if !g.started {
		panic("Can't tick a game that hasn't started!")
	}

	currentPlayer := g.Player1
	if g.CurrentTurn == P2 {
		currentPlayer = g.Player2
	}

	markIdx := currentPlayer.MakeMove(g.Board, g.CurrentTurn)
	if g.Board[markIdx] != EMPTY {
		panic("You're trying to mark an already marked cell!")
	}
	defer func() {
		if g.CurrentTurn == P1 {
			g.CurrentTurn = P2
		} else {
			g.CurrentTurn = P1
		}
	}()

	if g.CurrentTurn == P1 {
		g.Board[markIdx] = X
	} else {
		g.Board[markIdx] = O
	}
	g.filledCells += 1
	g.stats.MovementCount += 1

	patterns := [][]int{
		// Horizontal lines
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},

		// Vertical lines
		{0, 3, 6},
		{1, 4, 7},
		{2, 5, 8},

		// Diagonal lines
		{0, 4, 8},
		{6, 4, 2},
	}

patternLoop:
	for _, pattern := range patterns {
		firstCellMark := g.Board[pattern[0]]
		if firstCellMark == EMPTY {
			continue
		}

		for i := 1; i < 3; i++ {
			if g.Board[pattern[i]] != firstCellMark {
				continue patternLoop
			}
		}

		if g.CurrentTurn == P1 {
			g.stats.Winner = "Player 1"
		} else {
			g.stats.Winner = "Player 2"
		}
		g.stats.EndTime = time.Now()
		g.shouldEnd = true
	}

	if g.filledCells == 9 {
		g.stats.Winner = "No one"
		g.shouldEnd = true
		g.stats.EndTime = time.Now()
	}
}

func (g Game) Report() GameStats {
	return g.stats
}
