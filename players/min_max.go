package players

import (
	"math"
	"math/rand"

	sim "github.com/DanielRasho/IA-lab6/simulation"
)

type MinMaxPlayer struct{}

func NewMiniMaxPlayer() MinMaxPlayer {
	return MinMaxPlayer{}
}

func (player MinMaxPlayer) MakeMove(board sim.TicTacToeBoard, whoami sim.Turn) int {
	bestScore := math.MinInt
	// bestMove := -1
	bestMoves := make([]int, 0)

	playerMark, _ := sim.GetMarks(whoami)

	for _, move := range sim.GetAvailableCells(board) {
		board[move] = playerMark
		score := miniMax(board, sim.GetOpponent(whoami), true)
		board[move] = sim.EMPTY
		if score > bestScore {
			bestScore = score
			// bestMove = move
			bestMoves = nil
			bestMoves = make([]int, 0, 1)
			bestMoves = append(bestMoves, move)
		} else if score == bestScore {
			bestMoves = append(bestMoves, move)
		}
	}
	return bestMoves[rand.Intn(len(bestMoves))]
}

func miniMax(board sim.TicTacToeBoard, whoami sim.Turn, isMaximizing bool) int {
	score := evaluate(board, whoami)

	if score != -2 {
		return score
	}

	playerMark, oponentMark := sim.GetMarks(whoami)
	oponent := sim.GetOpponent(whoami)

	if isMaximizing {
		maxScore := math.MinInt
		for _, move := range sim.GetAvailableCells(board) {
			board[move] = playerMark
			score := miniMax(board, oponent, false)
			board[move] = sim.EMPTY
			maxScore = max(maxScore, score)
		}
		return maxScore
	} else {
		minScore := math.MaxInt
		for _, move := range sim.GetAvailableCells(board) {
			board[move] = oponentMark
			score := miniMax(board, whoami, true)
			board[move] = sim.EMPTY
			minScore = min(minScore, score)
		}
		return minScore
	}
}

// Evaluates the current node, it returns 10 if current player is the winner, -10,
// if it is the oponent, 0 tie, and -2 if there is not conclusion
func evaluate(board sim.TicTacToeBoard, whoami sim.Turn) int {
	// Get marks for both players
	playerMark, opponentMark := sim.GetMarks(whoami)

	// Define all possible winning lines
	winningLines := [][]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Columns
		{0, 4, 8}, {2, 4, 6}, // Diagonals
	}

	// Check for a win
	for _, line := range winningLines {
		if board[line[0]] == playerMark && board[line[1]] == playerMark && board[line[2]] == playerMark {
			return 10 // `whoami` wins
		}
		if board[line[0]] == opponentMark && board[line[1]] == opponentMark && board[line[2]] == opponentMark {
			return -10 // Opponent wins
		}
	}

	if len(sim.GetAvailableCells(board)) == 0 {
		return 0 // It's a tie
	}

	// Game is still in progress
	return -2
}
