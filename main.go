package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/DanielRasho/IA-lab6/players"
	"github.com/DanielRasho/IA-lab6/simulation"
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

type PlayerValue struct {
	val PlayerType
}

func (s PlayerValue) Set(obj string) error {
	switch obj {
	case "h":
		s.val = HUMAN
	case "m":
		s.val = MINIMAX
	case "a":
		s.val = ALFA_BETA
	case "mo":
		s.val = MONTECARLO
	default:
		return errors.New("Invalid player type provided!")
	}
	return nil
}

func (s PlayerValue) Get() interface{} {
	return s.val
}

func (s PlayerValue) String() string {
	switch s.val {
	case HUMAN:
		return "h"
	case MONTECARLO:
		return "mo"
	case MINIMAX:
		return "m"
	case ALFA_BETA:
		return "a"
	default:
		return ""
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

	pTypeFormat := "The type of player %d. Valid values are:\n* h: Human\n* m: Normal minimax\n* a: Minimax with alfa/beta prunning\n* mo: Montecarlo AI"
	var p1 PlayerValue
	flag.Var(&p1, "p1", fmt.Sprintf(pTypeFormat, 1))
	params.Player1 = p1.Get().(PlayerType)

	var p2 PlayerValue
	flag.Var(&p2, "p2", fmt.Sprintf(pTypeFormat, 2))
	params.Player2 = p2.Get().(PlayerType)

	flag.Parse()
	return params
}

func NewPlayerFromType(p PlayerType) simulation.TicTacToePlayer {
	switch p {
	case HUMAN:
		return players.NewHumanPlayer()
	default:
		panic("Player type not implemented!")
	}
}

func main() {
	params := ParseProgramParams()

	fmt.Printf("Constructing game of %s againts %s...\n", PrettyDisplay(params.Player1), PrettyDisplay(params.Player2))
	player1 := NewPlayerFromType(params.Player1)
	player2 := NewPlayerFromType(params.Player2)

	game := simulation.NewGameWith(player1, player2)

	fmt.Println("Starting game...")
	game.Start()
	for !game.ShouldEnd() {
		game.Tick()
	}
	fmt.Println(game.Report())
}
