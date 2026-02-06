package game

import (
	"os"
	"strconv"
	"strings"
)

func loadHighScore() int {
	data, err := os.ReadFile("highscore.dat")
	if err != nil {
		return 0
	}
	score, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0
	}
	return score
}

func saveHighScore(score int) {
	os.WriteFile("highscore.dat", []byte(strconv.Itoa(score)), 0644)
}
