package players

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	sim "github.com/DanielRasho/IA-lab6/simulation"
)

type HumanPlayer struct{ reader bufio.Reader }

func (self HumanPlayer) MakeMove(board sim.TicTacToeBoard) int {
	line, err := self.reader.ReadString('\n')
	if err != nil {
		panic("Failed to get user input!")
	}

	cellId, err := strconv.Atoi(strings.TrimSpace(line))
	for err != nil {
		cellId, err = strconv.Atoi(line)
	}

	return cellId - 1
}

func NewHumanPlayer() HumanPlayer {
	return HumanPlayer{reader: *bufio.NewReader(os.Stdin)}
}
