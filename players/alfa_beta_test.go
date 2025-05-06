package players

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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

func NodeToString[T any](node *AlfBetNode[T], b *strings.Builder) {
	b.WriteString("( v: ")
	b.WriteString(fmt.Sprintf("%v", node.Data))
	b.WriteString("; p: ")
	b.WriteString(strconv.FormatInt(int64(node.Point), 10))
	b.WriteString("; alfa: ")
	b.WriteString(strconv.FormatInt(int64(node.Alpha), 10))
	b.WriteString("; beta: ")
	b.WriteString(strconv.FormatInt(int64(node.Beta), 10))
	b.WriteString(" )")
}

func _TreeToString[T any](root *AlfBetNode[T], b *strings.Builder, prefix string) {
	if root == nil { // Just in case
		return
	}

	b.WriteString(prefix)
	NodeToString(root, b)
	b.WriteRune('\n')
	for _, v := range root.Children {
		_TreeToString(v, b, prefix+"  * ")
	}
}

func TreeToString[T any](root *AlfBetNode[T]) string {
	b := strings.Builder{}
	_TreeToString(root, &b, "")
	return b.String()
}

func _compareTrees[T comparable](expected *AlfBetNode[T], actual *AlfBetNode[T]) bool {
	if !expected.EqualByMetadata(actual) {
		return false
	}

	for _, exChild := range expected.Children {
		foundChild := false
		childIdx := 0
		for i, acChild := range actual.Children {
			if exChild.EqualByMetadata(acChild) {
				foundChild = true
				childIdx = i
				break
			}
		}

		if !foundChild {
			return false
		}

		subTreeIsEqual := _compareTrees(exChild, actual.Children[childIdx])
		if !subTreeIsEqual {
			return false
		}
	}

	return true
}

func compareTrees[T comparable](t *testing.T, expected *AlfBetNode[T], actual *AlfBetNode[T]) {
	areEqual := _compareTrees(expected, actual)
	if !areEqual {
		t.Logf("\nEXPECTED:\n%s", TreeToString(expected))
		t.Logf("\nACTUAL:\n%s", TreeToString(actual))
		t.Fatalf("Actual tree != Expected Tree!")
	}
}

func TestAlfaBetaPrunning(t *testing.T) {
	board := TwoEmptyCellsBoard()
	player := AlfBetPlayer{}
	player.CreateTreeFromBoard(&board, sim.P2, true)

	expectedTree := CreateTwoEmptyCellsExpectedTree()
	compareTrees(t, expectedTree, player.Tree)
}
