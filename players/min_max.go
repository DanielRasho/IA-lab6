package players

import (
	"math"

	sim "github.com/DanielRasho/IA-lab6/simulation"
)

type MinMaxPlayer struct{}

func NewMiniMaxPlayer() MinMaxPlayer {
	return MinMaxPlayer{}
}

func (self MinMaxPlayer) MakeMove(board sim.TicTacToeBoard, whoami sim.Turn) int {
	bestScore := math.MinInt
	bestMove := -1

	playerMark, _ := sim.GetMarks(whoami)

	for _, move := range sim.GetAvailableCells(board) {
		board[move] = playerMark
		score := miniMax(board, sim.GetOpponent(whoami), false, whoami)
		board[move] = sim.EMPTY
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}
	return bestMove
}

func miniMax(board sim.TicTacToeBoard, turn sim.Turn, isMaximizing bool, whoami sim.Turn) int {
	score := evaluate(board, whoami)

	if score != -2 {
		return score
	}

	playerMark, oponentMark := sim.GetMarks(turn)
	oponent := sim.GetOpponent(turn)

	if isMaximizing {
		maxScore := math.MinInt
		for _, move := range sim.GetAvailableCells(board) {
			board[move] = playerMark
			score := miniMax(board, oponent, false, whoami)
			board[move] = sim.EMPTY
			maxScore = max(maxScore, score)
		}
		return maxScore
	} else {
		minScore := math.MaxInt
		for _, move := range sim.GetAvailableCells(board) {
			board[move] = oponentMark
			score := miniMax(board, turn, true, whoami)
			board[move] = sim.EMPTY
			minScore = min(minScore, score)
		}
		return minScore
	}
}

func evaluate(board sim.TicTacToeBoard, whoami sim.Turn) int {
	winner := sim.GetBoardWinner(sim.ToBitMasks(&board))
	if winner == nil {
		// No winner, score -2, keep searching
		return -2
	} else {
		if *winner == whoami {
			return 10
		} else if *winner == sim.GetOpponent(whoami) {
			return -10
		}
	}
	// It's a draw if all cells are filled and no winner
	if len(sim.GetAvailableCells(board)) == 0 {
		return 0
	}
	return -2
}
