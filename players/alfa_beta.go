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
	Point float64
	Alpha float64
	Beta  float64
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
func GetPoint[T any](self *AlfBetNode[T]) float64 {
	return self.Point
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

func approxEqual(a float64, b float64) bool {
	epsilon := math.Nextafter(1.0, 2.0) - 1.0
	return bothAreNotInf(a, b) && a-b <= epsilon
}

// a<=b with floats64
func lessOrAproxEqual(a float64, b float64) bool {
	return a < b || approxEqual(a, b)
}

func isInf(a float64) bool {
	return math.IsInf(a, 1) || math.IsInf(a, -1)
}

func bothAreNotInf(a float64, b float64) bool {
	return !isInf(a) && !isInf(b)
}

func mapAndApply[T any](data []T, extractor func(T) float64, op func(float64, float64) float64) float64 {

	total := 0.0
	for _, v := range data {
		total = op(total, extractor(v))
	}

	return total
}

func max2(a, b float64) float64 {
	return max(a, b)
}

func min2(a, b float64) float64 {
	return min(a, b)
}

func _createTreeFromBoard(root *AlfBetNode[AlfBetNodeData], board *sim.TicTacToeBoard, original sim.Turn, whoami sim.Turn, treeDepth int) (float64, float64, bool) {
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
				root.Point = 10 / float64(treeDepth)
			case sim.GetOpponent(original):
				// Opponent won!
				root.Point = -10 / float64(treeDepth)
			}
		}

	} else {
		for i, v := range *board {
			if v == sim.EMPTY {
				markedMask := sim.CopyAndMark(boardMask, currentMark, i)
				markedBoard := sim.BoardFromBitMask(markedMask)

				child := &AlfBetNode[AlfBetNodeData]{
					Data:   AlfBetNodeData{BoardBitMask: markedMask},
					Parent: root,
					Alpha:  root.Alpha,
					Beta:   root.Beta,
					IsMax:  !root.IsMax,
				}
				root.Children = append(root.Children, child)

				alpha, beta, wasChildNodeLeaf := _createTreeFromBoard(child, &markedBoard, original, currentOponent, treeDepth+1)

				if root.IsMax {
					if child.Point > root.Alpha {
						root.Data.Move = i
					}
					alpha = max(alpha, root.Alpha, child.Point)
				} else {
					if child.Point < root.Beta {
						root.Data.Move = i
					}
					beta = min(beta, root.Beta, child.Point)
				}

				if wasChildNodeLeaf {
					root.Alpha, root.Beta = max(alpha, root.Alpha), min(beta, root.Beta)
				} else {
					if !approxEqual(root.Alpha, child.Alpha) && !isInf(child.Alpha) {
						root.Beta = min(root.Beta, child.Alpha)
					} else if !approxEqual(root.Beta, child.Beta) && !isInf(child.Beta) {
						root.Alpha = max(root.Alpha, child.Beta)
					}
				}

				if root.IsMax {
					root.Point = mapAndApply(root.Children, GetPoint, max2)
				} else {
					root.Point = mapAndApply(root.Children, GetPoint, min2)
				}

				if lessOrAproxEqual(root.Beta, root.Alpha) {
					break
				}
			}
		}
	}

	return root.Alpha, root.Beta, isLeafNode
}

func (self *AlfBetPlayer) CreateTreeFromBoard(board *sim.TicTacToeBoard, whoami sim.Turn, isMax bool) {
	self.Tree = &AlfBetNode[AlfBetNodeData]{
		Alpha: math.Inf(-1),
		Beta:  math.Inf(1),
		IsMax: isMax,
		Data:  AlfBetNodeData{BoardBitMask: sim.ToBitMasks(board)},
	}

	_createTreeFromBoard(self.Tree, board, whoami, whoami, 1)
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
		// s.CreateTreeFromBoard(&board, whoami, whoami == sim.P1)
		s.CreateTreeFromBoard(&board, whoami, true)
		move = s.FindMoveOnTree(boardBitMask)
	}

	if move == nil {
		panic("Move not found even after regenerating tree!")
	}

	return *move
}
