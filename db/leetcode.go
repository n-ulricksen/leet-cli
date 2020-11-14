package db

import (
	"fmt"

	"github.com/ulricksennick/lcfetch/problem"
)

const (
	leetcodeApiUrl = "https://leetcode.com/api/problems/all/"
)

func main() {
	// TODO: write proper tests and delete this file

	// dropProblems()

	// fetchAndParseProblems()

	// printAllProblems()

	// testFilters()

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// Test stuff
func dropProblems() {
	database, err := CreateDB()
	must(err)
	database.DropAllProblems()
}

func testFilters() {
	database, err := CreateDB()
	must(err)
	probs, err := database.GetAllProblems()

	difficulty := problem.HARD

	filtered := problem.FilterByDifficulty(probs, difficulty)
	filtered = problem.FilterOutPaid(filtered)
	for _, prob := range filtered {
		fmt.Printf("%v\t%v\n", prob.DisplayId, prob.Name)
	}
	fmt.Println()
	fmt.Printf("%d filtered problems.\n", len(filtered))
}

func printAllProblems() {
	database, err := CreateDB()
	must(err)
	probs, err := database.GetAllProblems()
	must(err)
	for _, prob := range probs {
		fmt.Printf("%v\t%v\n", prob.DisplayId, prob.Name)
	}
	fmt.Printf("Fetched %d problems.\n", len(probs))
}
