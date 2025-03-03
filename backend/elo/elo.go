package elo

import "math"

const (
	KFactorHigh = 32
	KFactorMedium = 24
	KFactorLow = 16
)

func CalculateElo(winnerElo, loserElo int, winner bool) (int, int) {
	// probability winning for each
	winnerProb := 1 / (1 + math.Pow(10, float64(loserElo - winnerElo)/400))
	loserProb := 1 / (1 + math.Pow(10, float64(winnerElo - loserElo)/400))

	winnerK := getKFactor(winnerElo)
	loserK := getKFactor(loserElo)

	newWinnerElo := winnerElo + int(float64(winnerK) * (1 - winnerProb))
	newLoserElo := loserElo + int(float64(loserK) * (0 - loserProb))

	return newWinnerElo, newLoserElo
}

func getKFactor(elo int) int {
	switch {
	case elo < 1600:
		return KFactorHigh
	case elo < 2000:
		return KFactorMedium
	default:
		return KFactorLow
	}
}