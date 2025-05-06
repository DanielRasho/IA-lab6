package players

import (
	"math"
	"testing"

	sim "github.com/DanielRasho/IA-lab6/simulation"
)

func NewNode[T any](alpha int, beta int, parent *AlfBetNode[T]) *AlfBetNode[T] {

	node := &AlfBetNode[T]{
		Alpha:  alpha,
		Beta:   beta,
		Parent: parent,
	}

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node
}

func NewIntNode(alpha int, beta int, parent *AlfBetNode[int]) *AlfBetNode[int] {
	return NewNode(alpha, beta, parent)
}

func NewChildFrom[T any](parent *AlfBetNode[T]) *AlfBetNode[T] {
	return NewNode(0, 0, parent)
}

func TwoEmptyCellsBoard() sim.TicTacToeBoard {
	return []sim.CellMark{
		sim.X, sim.O, sim.X,
		sim.X, sim.O, sim.X,
		sim.EMPTY, sim.EMPTY, sim.O,
	}
}

func CreateTwoEmptyCellsExpectedTree() *AlfBetNode[int] {
	root := NewIntNode(math.MinInt, math.MaxInt, nil)
	root.Alpha = 0
	root.Beta = 1

	childLeft := NewChildFrom(root)
	childLeft.Alpha = math.MinInt
	childLeft.Beta = 0

	childRight := NewChildFrom(root)
	childRight.Alpha = 1
	childLeft.Beta = math.MaxInt
	childRight.Point = 1

	grandChildLeft := NewChildFrom(childLeft)
	grandChildLeft.Point = 0

	return root
}

func compareTrees[T comparable](t *testing.T, expected *AlfBetNode[T], actual *AlfBetNode[T]) {

}

func TestAlfaBetaPrunning(t *testing.T) {
	board := TwoEmptyCellsBoard()
	player := AlfBetPlayer{}
	player.CreateTreeFromBoard(&board, sim.P2, true)

	expectedTree := CreateTwoEmptyCellsExpectedTree()
	compareTrees(t, expectedTree, player.Tree)
}
