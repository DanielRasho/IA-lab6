package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/DanielRasho/IA-lab6/players"
	sim "github.com/DanielRasho/IA-lab6/simulation"
)

type PlayerType = int

const (
	HUMAN PlayerType = iota
	MINIMAX
	ALFA_BETA
	MONTECARLO
)

func PrettyDisplay(p PlayerType) string {
	switch p {
	case HUMAN:
		return "Human"
	case MINIMAX:
		return "Minimax"
	case ALFA_BETA:
		return "Alfa/Beta prunning"
	case MONTECARLO:
		return "Montecarlo"
	default:
		panic("Undefined player type!")
	}
}

func FromString(s string) PlayerType {
	switch strings.TrimSpace(s) {
	case "h":
		return HUMAN
	case "mi":
		return MINIMAX
	case "a":
		return ALFA_BETA
	case "mo":
		return MONTECARLO
	default:
		panic("Invalid player type supplied!")
	}
}

type ProgramParams struct {
	ShowBoard bool
	Player1   PlayerType
	Player2   PlayerType
}

func ParseProgramParams() ProgramParams {
	params := ProgramParams{}
	flag.BoolVar(&params.ShowBoard, "showBoard", true, "Describes whether or not the board should be displayed when making a game")

	pTypeFormat := "The type of player %d. Valid values are:\n* h: Human\n* mi: Normal minimax\n* a: Minimax with alfa/beta prunning\n* mo: Montecarlo AI"

	var p1, p2 string
	flag.StringVar(&p1, "p1", "h", fmt.Sprintf(pTypeFormat, 1))
	flag.StringVar(&p2, "p2", "h", fmt.Sprintf(pTypeFormat, 2))

	flag.Parse()
	params.Player1 = FromString(p1)
	params.Player2 = FromString(p2)
	return params
}

func NewPlayerFromType(p PlayerType) sim.TicTacToePlayer {
	switch p {
	case HUMAN:
		return players.NewHumanPlayer()
	case MINIMAX:
		return players.NewMiniMaxPlayer()
	case ALFA_BETA:
		return players.NewAlfaBetaPlayer()
	default:
		panic("Player type not implemented!")
	}
}

func main() {
	params := ParseProgramParams()

	fmt.Printf("Constructing game of %s againts %s...\n", PrettyDisplay(params.Player1), PrettyDisplay(params.Player2))
	player1 := NewPlayerFromType(params.Player1)
	player2 := NewPlayerFromType(params.Player2)

	game := sim.NewGameWith(player1, player2)

	fmt.Println("Starting game...")
	game.Start()
	// game.Board = []sim.CellMark{
	// 	sim.X, sim.O, sim.X,
	// 	sim.O, sim.EMPTY, sim.X,
	// 	sim.EMPTY, sim.EMPTY, sim.O,
	// }
	// game.CurrentTurn = sim.P2
	for !game.ShouldEnd() {
		if params.ShowBoard {
			DisplayBoard(game)
			fmt.Println("")
		}
		game.Tick()
	}

	if params.ShowBoard {
		DisplayBoard(game)
		fmt.Println("")
	}
	DisplayReport(game.Report())
}
