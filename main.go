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
	ShowBoard      bool
	Player1        PlayerType
	Player2        PlayerType
	RunExperiments bool
}

func RunExperiments() {
	fmt.Println("Running Monte Carlo Tree Search experiments...")
	
	wins := 0
	draws := 0
	losses := 0
	totalNodesExplored := 0
	
	mctsPlayer := players.NewMonteCarloPlayer()
	minimaxPlayer := players.NewMiniMaxPlayer()
	
	N := 1000
	for i := 0; i < N; i++ {
		var game sim.Game
		var mctsPlayerTurn sim.Turn
		
		if i%2 == 0 {
			game = sim.NewGameWith(mctsPlayer, minimaxPlayer)
			mctsPlayerTurn = sim.P1
		} else {
			game = sim.NewGameWith(minimaxPlayer, mctsPlayer)
			mctsPlayerTurn = sim.P2
		}
		
		game.Start()
		
		nodesExplored := 0
		
		for !game.ShouldEnd() {
			if game.CurrentTurn == mctsPlayerTurn {
				nodesExplored++
			}
			
			game.Tick()
		}
		
		report := game.Report()
		if report.Winner == "No one" {
			draws++
		} else if (mctsPlayerTurn == sim.P1 && report.Winner == "Player 1") ||
			(mctsPlayerTurn == sim.P2 && report.Winner == "Player 2") {
			wins++
		} else {
			losses++
		}
		
		totalNodesExplored += nodesExplored
	}
	
	fmt.Printf("Results after %d experiments:\n", N)
	fmt.Printf("Wins: %d (%.2f%%)\n", wins, float64(wins)/float64(N)*100)
	fmt.Printf("Draws: %d (%.2f%%)\n", draws, float64(draws)/float64(N)*100)
	fmt.Printf("Losses: %d (%.2f%%)\n", losses, float64(losses)/float64(N)*100)
	fmt.Printf("Average nodes explored per game: %.2f\n", float64(totalNodesExplored)/float64(N))
}

func ParseProgramParams() ProgramParams {
	params := ProgramParams{}
	flag.BoolVar(&params.ShowBoard, "showBoard", true, "Describes whether or not the board should be displayed when making a game")
	flag.BoolVar(&params.RunExperiments, "experiments", false, "Run MCTS experiments")

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
	case MONTECARLO:
		return players.NewMonteCarloPlayer()
	default:
		panic("Player type not implemented!")
	}
}


func main() {
	params := ParseProgramParams()

	if params.RunExperiments {
		RunExperiments()
		return
	}

	fmt.Printf("Constructing game of %s againts %s...\n", PrettyDisplay(params.Player1), PrettyDisplay(params.Player2))
	player1 := NewPlayerFromType(params.Player1)
	player2 := NewPlayerFromType(params.Player2)

	game := sim.NewGameWith(player1, player2)

	fmt.Println("Starting game...")
	game.Start()
	// | 1 | 2 | O |
	// | X | X | O |
	// | 7 | 8 | X |
	// game.Board = []sim.CellMark{
	// 	sim.EMPTY, sim.EMPTY, sim.O,
	// 	sim.X, sim.X, sim.O,
	// 	sim.EMPTY, sim.EMPTY, sim.X,
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
