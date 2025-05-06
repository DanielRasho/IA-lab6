package players

import (
	"math"

	sim "github.com/DanielRasho/IA-lab6/simulation"
)

type AlfBetNode[T any] struct {
	Data     T
	Children []*AlfBetNode[T]
	Parent   *AlfBetNode[T]

	IsMax bool
	Point int
	Alpha int
	Beta  int
}

type AlfBetPlayer struct {
	Tree *AlfBetNode[int]
}

func countBits(mask int) int {
	count := 0
	for mask > 0 {
		mask &= (mask - 1)
		count++
	}
	return count
}

func _createTreeFromBoard[T any](root *AlfBetNode[T], board *sim.TicTacToeBoard, original sim.Turn, whoami sim.Turn) (int, int) {
	currentMark, _ := sim.GetMarks(whoami)
	currentOponent := sim.GetOpponent(whoami)

	boardMask := sim.ToBitMasks(board)
	isLeafNode := countBits(boardMask) == 9
	if isLeafNode {
		winner := sim.GetBoardWinner(boardMask)
		if winner == nil {
			// No winner
			// Leaf node has a value of 0
			root.Point = 0
		}

		winnerVal := *winner
		switch winnerVal {
		case original:
			// I won!
			// Leaf node has a value of 1
			root.Point = 1
		case sim.GetOpponent(original):
			// Opponent won!
			// Leaf node has a value of -1
			root.Point = -1
		}

	} else {
		for i, v := range *board {
			if v == sim.EMPTY {
				markedMask := sim.CopyAndMark(boardMask, currentMark, i)
				markedBoard := sim.BoardFromBitMask(markedMask)

				child := AlfBetNode[T]{
					Parent: root,
					Alpha:  root.Alpha,
					Beta:   root.Beta,
				}

				alpha, beta := _createTreeFromBoard(&child, &markedBoard, original, currentOponent)
				if root.IsMax {
					alpha = max(alpha, root.Alpha, child.Point)
				} else {
					beta = min(beta, root.Beta, child.Point)
				}

				if root.Alpha != alpha {
					root.Beta = alpha
				} else if root.Beta != beta {
					root.Alpha = beta
				}
			}
		}
	}

	return root.Alpha, root.Beta
}

func (self *AlfBetPlayer) CreateTreeFromBoard(board *sim.TicTacToeBoard, whoami sim.Turn, isMax bool) {
	self.Tree = &AlfBetNode[int]{
		Alpha: math.MinInt,
		Beta:  math.MaxInt,
		IsMax: isMax,
	}

	_createTreeFromBoard(self.Tree, board, whoami, whoami)
}
