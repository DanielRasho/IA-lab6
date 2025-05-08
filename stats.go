package main

import (
	"sync"
)

type PlayerPair struct {
	P1 PlayerType
	P2 PlayerType
}

type StatsRow struct {
	WinsP1    float64
	Draws     float64
	WinsP2    float64
	GameCount float64
}

type ReportTable = map[PlayerPair]StatsRow

type GameEndedStats struct {
	P1     PlayerType
	P2     PlayerType
	Winner *PlayerType
}

func ListenForStats(wg *sync.WaitGroup, gameEnded chan GameEndedStats, table *ReportTable) {
	defer wg.Done()

	for gameResult := range gameEnded {
		key := PlayerPair{
			P1: gameResult.P1,
			P2: gameResult.P2,
		}
		val := *table

		if _, found := val[key]; !found {
			val[key] = StatsRow{}
		}

		previous := val[key]
		previous.GameCount += 1

		if gameResult.Winner == nil {
			previous.Draws += 1
		} else if *gameResult.Winner == gameResult.P1 {
			previous.WinsP1 += 1
		} else if *gameResult.Winner == gameResult.P2 {
			previous.WinsP2 += 1
		}

		val[key] = previous
	}
}
