package players

import (
	"math"

	"github.com/DanielRasho/IA-lab6/simulation"
)

type AlfBetNode[T any] struct {
	Data     T
	Children []AlfBetNode[T]

	alpha int
	beta  int
}

type AlfBetPlayer struct {
	tree *AlfBetNode[int]
}

func _createTreeFromBoard[T any](root *AlfBetNode[T], board *simulation.TicTacToeBoard) {

}

func (self *AlfBetPlayer) CreateTreeFromBoard(board *simulation.TicTacToeBoard) {

	root := AlfBetNode[int]{
		alpha: math.MaxInt,
		beta:  math.MinInt,
	}

	_createTreeFromBoard(&root, board)
}
