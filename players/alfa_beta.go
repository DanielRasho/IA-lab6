package players

import (
	"math"

	sim "github.com/DanielRasho/IA-lab6/simulation"
)

type AlfBetNodeData struct {
	// The state the board should be at this point.
	BoardBitMask int
	// 0-8 index that the bot should mark when it encounters this state.
	Move int
}

type AlfBetNode[T any] struct {
	Data     T
	Children []*AlfBetNode[T]
	Parent   *AlfBetNode[T]

	IsMax bool
	Point int
	Alpha int
	Beta  int
}

func EqualByData[T comparable](self *AlfBetNode[T], other *AlfBetNode[T]) bool {
	return self.Data == other.Data
}
func (self *AlfBetNode[T]) EqualByMetadata(other *AlfBetNode[T]) bool {
	return self.Point == other.Point &&
		self.Alpha == other.Alpha &&
		self.Beta == other.Beta &&
		self.IsMax == other.IsMax
}

type AlfBetPlayer struct {
	Tree *AlfBetNode[AlfBetNodeData]
}

func countBits(mask int) int {
	count := 0
	for mask > 0 {
		mask &= (mask - 1)
		count++
	}
	return count
}

func _createTreeFromBoard[T any](root *AlfBetNode[T], board *sim.TicTacToeBoard, original sim.Turn, whoami sim.Turn) (int, int, bool) {
	currentMark, _ := sim.GetMarks(whoami)
	currentOponent := sim.GetOpponent(whoami)

	boardMask := sim.ToBitMasks(board)
	winner := sim.GetBoardWinner(boardMask)
	isLeafNode := countBits(boardMask) == 9 || winner != nil
	if isLeafNode {
		if winner == nil {
			// No winner
			// Leaf node has a value of 0
			root.Point = 0
		} else {
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
		}

	} else {
		for i, v := range *board {
			if v == sim.EMPTY {
				markedMask := sim.CopyAndMark(boardMask, currentMark, i)
				markedBoard := sim.BoardFromBitMask(markedMask)

				child := &AlfBetNode[T]{
					Parent: root,
					Alpha:  root.Alpha,
					Beta:   root.Beta,
					IsMax:  !root.IsMax,
				}
				root.Children = append(root.Children, child)

				alpha, beta, wasChildNodeLeaf := _createTreeFromBoard(child, &markedBoard, original, currentOponent)
				if root.IsMax {
					alpha = max(alpha, root.Alpha, child.Point)
				} else {
					beta = min(beta, root.Beta, child.Point)
				}

				if !wasChildNodeLeaf {
					if root.Alpha != alpha {
						root.Beta = alpha
					} else if root.Beta != beta {
						root.Alpha = beta
					}
				} else {
					root.Alpha, root.Beta = alpha, beta
				}

				if root.IsMax {
					root.Point = root.Alpha
				} else {
					root.Point = root.Beta
				}
			}
		}
	}

	return root.Alpha, root.Beta, isLeafNode
}

func (self *AlfBetPlayer) CreateTreeFromBoard(board *sim.TicTacToeBoard, whoami sim.Turn, isMax bool) {
	self.Tree = &AlfBetNode[AlfBetNodeData]{
		Alpha: math.MinInt,
		Beta:  math.MaxInt,
		IsMax: isMax,
	}

	_createTreeFromBoard(self.Tree, board, whoami, whoami)
}

func _FindMoveOnTree(mask int, current *AlfBetNode[AlfBetNodeData]) *int {
	if current == nil {
		return nil
	}

	if current.Data.BoardBitMask == mask {
		return &current.Data.Move
	}

	for _, v := range current.Children {
		foundMove := _FindMoveOnTree(mask, v)
		if foundMove != nil {
			return foundMove
		}
	}

	return nil
}

// Attempts to find a move on the tree.
// If no tree exists or the move couldn't be found the nil is returned.
func (self *AlfBetPlayer) FindMoveOnTree(boardBitMask int) *int {
	if self.Tree == nil {
		return nil
	}

	return _FindMoveOnTree(boardBitMask, self.Tree)
}

func NewAlfaBetaPlayer() AlfBetPlayer {
	return AlfBetPlayer{}
}

func (s AlfBetPlayer) MakeMove(board sim.TicTacToeBoard, whoami sim.Turn) int {
	// If we're player 1 we always start in the middle
	boardBitMask := sim.ToBitMasks(&board)
	if countBits(boardBitMask) == 0 && whoami == sim.P1 {
		return 4
	}

	// If we're player 2 or second turn of player 1, calculate tree only if nil

	move := s.FindMoveOnTree(boardBitMask)
	if move == nil {
		s.CreateTreeFromBoard(&board, whoami, true)
		move = s.FindMoveOnTree(boardBitMask)
	}

	if move == nil {
		panic("Move not found even after regenerating tree!")
	}

	return *move
}
