package players

import (
	"math"
	"math/rand"
	"time"

	sim "github.com/DanielRasho/IA-lab6/simulation"
)

type MCTSNode struct {
	State       sim.TicTacToeBoard
	Parent      *MCTSNode
	Children    []*MCTSNode
	Visits      int
	Wins        int
	Draws       int
	Losses      int
	UntriedMoves []int
	PlayerTurn  sim.Turn
	Move        int
}

func NewMCTSNode(state sim.TicTacToeBoard, parent *MCTSNode, move int, playerTurn sim.Turn) *MCTSNode {
	node := &MCTSNode{
		State:       make(sim.TicTacToeBoard, len(state)),
		Parent:      parent,
		Children:    []*MCTSNode{},
		Visits:      0,
		Wins:        0,
		Draws:       0,
		Losses:      0,
		PlayerTurn:  playerTurn,
		Move:        move,
	}
	
	copy(node.State, state)
	
	node.UntriedMoves = sim.GetAvailableCells(state)
	
	return node
}

func (node *MCTSNode) UCB1(explorationWeight float64) float64 {
	if node.Visits == 0 {
		return math.Inf(1)
	}
	
	exploitation := float64(node.Wins) / float64(node.Visits)
	exploration := explorationWeight * math.Sqrt(math.Log(float64(node.Parent.Visits))/float64(node.Visits))
	
	return exploitation + exploration
}

func (node *MCTSNode) SelectChild() *MCTSNode {
	bestScore := math.Inf(-1)
	var bestChild *MCTSNode
	
	for _, child := range node.Children {
		score := child.UCB1(1.414)
		if score > bestScore {
			bestScore = score
			bestChild = child
		}
	}
	
	return bestChild
}

func (node *MCTSNode) Expand() *MCTSNode {
	if len(node.UntriedMoves) == 0 {
		return nil
	}
	
	moveIndex := rand.Intn(len(node.UntriedMoves))
	move := node.UntriedMoves[moveIndex]
	
	node.UntriedMoves = append(node.UntriedMoves[:moveIndex], node.UntriedMoves[moveIndex+1:]...)
	
	newState := make(sim.TicTacToeBoard, len(node.State))
	copy(newState, node.State)
	
	playerMark, _ := sim.GetMarks(node.PlayerTurn)
	newState[move] = playerMark
	
	childNode := NewMCTSNode(
		newState,
		node,
		move,
		sim.GetOpponent(node.PlayerTurn),
	)
	
	node.Children = append(node.Children, childNode)
	return childNode
}

func (node *MCTSNode) Simulate() int {
	tempState := make(sim.TicTacToeBoard, len(node.State))
	copy(tempState, node.State)
	
	currentTurn := node.PlayerTurn
	
	for {
		if sim.PlayerWonInBoard(&tempState, sim.GetOpponent(currentTurn)) {
			if sim.GetOpponent(currentTurn) == node.PlayerTurn {
				return 1
			} else {
				return -1
			}
		}
		
		availableMoves := sim.GetAvailableCells(tempState)
		if len(availableMoves) == 0 {
			return 0
		}
		
		moveIndex := rand.Intn(len(availableMoves))
		move := availableMoves[moveIndex]
		
		playerMark, _ := sim.GetMarks(currentTurn)
		tempState[move] = playerMark
		
		currentTurn = sim.GetOpponent(currentTurn)
	}
}

func (node *MCTSNode) Backpropagate(result int) {
	node.Visits++
	
	if result > 0 {
		node.Wins++
	} else if result < 0 {
		node.Losses++
	} else {
		node.Draws++
	}
	
	if node.Parent != nil {
		node.Parent.Backpropagate(-result)
	}
}

func (node *MCTSNode) BestMove() int {
	bestVisits := -1
	var bestMove int
	
	for _, child := range node.Children {
		if child.Visits > bestVisits {
			bestVisits = child.Visits
			bestMove = child.Move
		}
	}
	
	return bestMove
}

type MonteCarloPlayer struct {
	Iterations int
	TimeLimit  time.Duration
}

func NewMonteCarloPlayer() MonteCarloPlayer {
	return MonteCarloPlayer{
		Iterations: 10000,
		TimeLimit:  2 * time.Second,
	}
}

func (player MonteCarloPlayer) MakeMove(board sim.TicTacToeBoard, whoami sim.Turn) int {
	if countEmptyCells(board) == 9 {
		return 4
	}
	
	rootNode := NewMCTSNode(board, nil, -1, whoami)
	
	startTime := time.Now()
	iterations := 0
	
	for time.Since(startTime) < player.TimeLimit && iterations < player.Iterations {
		node := rootNode
		for len(node.UntriedMoves) == 0 && len(node.Children) > 0 {
			node = node.SelectChild()
		}
		
		if len(node.UntriedMoves) > 0 {
			node = node.Expand()
		}
		
		result := node.Simulate()
		
		node.Backpropagate(result)
		
		iterations++
	}
	
	return rootNode.BestMove()
}

func countEmptyCells(board sim.TicTacToeBoard) int {
	count := 0
	for _, cell := range board {
		if cell == sim.EMPTY {
			count++
		}
	}
	return count
}

