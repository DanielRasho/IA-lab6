package simulation

import (
	"fmt"
	"strings"
	"testing"
)

func BoardsSideBySideString(expected TicTacToeBoard, actual TicTacToeBoard) string {
	length := max(len(expected), len(actual))

	b := strings.Builder{}
	b.WriteString("\n")
	b.WriteString("  Expected:\t   Actual:\n")

	for i := 0; i < length; i += 3 {
		i2 := i + 1
		i3 := i + 2

		shouldPrintI := i < len(expected)
		if shouldPrintI {
			b.WriteString(fmt.Sprintf("| %s ", CellToString(expected[i], int64(i))))
		}
		shouldPrintI2 := i2 < len(expected)
		if shouldPrintI2 {
			b.WriteString(fmt.Sprintf("| %s ", CellToString(expected[i2], int64(i2))))
		}
		shouldPrintI3 := i3 < len(expected)
		if shouldPrintI3 {
			b.WriteString(fmt.Sprintf("| %s ", CellToString(expected[i3], int64(i3))))
		}

		printedExpected := shouldPrintI || shouldPrintI2 || shouldPrintI3
		if printedExpected {
			b.WriteString("|\t")
		} else {
			b.WriteString("\t\t")
		}

		shouldPrintAI := i < len(actual)
		if shouldPrintI {
			b.WriteString(fmt.Sprintf("| %s ", CellToString(actual[i], int64(i))))
		}
		shouldPrintAI2 := i2 < len(actual)
		if shouldPrintI2 {
			b.WriteString(fmt.Sprintf("| %s ", CellToString(actual[i2], int64(i2))))
		}
		shouldPrintAI3 := i3 < len(actual)
		if shouldPrintI3 {
			b.WriteString(fmt.Sprintf("| %s ", CellToString(actual[i3], int64(i3))))
		}

		shouldMakeNewLine := (shouldPrintI && shouldPrintI2 && shouldPrintI3) || (shouldPrintAI && shouldPrintAI2 && shouldPrintAI3)
		printedExpected = shouldPrintI || shouldPrintI2 || shouldPrintI3
		if printedExpected && shouldMakeNewLine {
			b.WriteString("|\n")
		} else if shouldMakeNewLine {
			b.WriteString("\n")
		}

	}

	return b.String()
}

func compareBoards(t *testing.T, expected TicTacToeBoard, actual TicTacToeBoard) {
	if len(expected) != len(actual) {
		t.Log(BoardsSideBySideString(expected, actual))
		t.Fatalf("Board lengths do not match! %d != %d", len(expected), len(actual))
	}

	for i, v := range expected {
		if v != actual[i] {
			t.Log(BoardsSideBySideString(expected, actual))
			t.Fatalf("Cell %d doesn't equal expected! %s != %s",
				i,
				CellToString(v, int64(i)),
				CellToString(actual[i], int64(i)))
		}
	}
}

func TestSimpleBoardConversion(t *testing.T) {
	board := []CellMark{
		X, EMPTY, EMPTY,
		EMPTY, EMPTY, EMPTY,
		EMPTY, EMPTY, EMPTY,
	}
	bitmask := ToBitMasks(&board)
	reparsedBoard := BoardFromBitMask(bitmask)

	expectedBitMask := 1

	if bitmask != expectedBitMask {
		t.Fatalf("Bitmask are not equal! %d != %d", expectedBitMask, bitmask)
	}

	compareBoards(t, board, reparsedBoard)
}

func TestComplexBoardConversion(t *testing.T) {
	board := []CellMark{
		X, O, X,
		X, O, X,
		EMPTY, EMPTY, O,
	}
	bitmask := ToBitMasks(&board)
	reparsedBoard := BoardFromBitMask(bitmask)

	compareBoards(t, board, reparsedBoard)
}

func TestBoardMark(t *testing.T) {
	initialBoard := []CellMark{
		X, EMPTY, EMPTY,
		EMPTY, EMPTY, EMPTY,
		EMPTY, EMPTY, EMPTY,
	}
	expectedBoard := []CellMark{
		X, EMPTY, EMPTY,
		EMPTY, EMPTY, EMPTY,
		EMPTY, EMPTY, O,
	}

	bitmask := ToBitMasks(&initialBoard)
	bitmask = CopyAndMark(bitmask, O, 8)
	board := BoardFromBitMask(bitmask)

	compareBoards(t, expectedBoard, board)

}
