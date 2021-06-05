package util

import "fmt"

// 8-bit/256-color lookup table found at:
// https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit
var colors = map[int]int{
	1: 82,
	2: 208,
	3: 196,
}

func GetColoredDifficultyText(difficulty int) string {
	var diff string

	switch difficulty {
	case 1:
		diff = "Easy"
	case 2:
		diff = "Medium"
	case 3:
		diff = "Hard"
	}
	diff = fmt.Sprintf("\033[38;5;%dm%s\033[m",
		colors[difficulty], diff)

	return diff
}
