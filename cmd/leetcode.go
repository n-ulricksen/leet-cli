package main

import (
	"fmt"
	"net/http"

	"github.com/ulricksennick/leetcode-fetcher/db"
	"github.com/ulricksennick/leetcode-fetcher/parser"
	"github.com/ulricksennick/leetcode-fetcher/problem"
)

const (
	leetcodeApiUrl = "https://leetcode.com/api/problems/all/"
)

func main() {
	fmt.Println("vim-go")

	// TODO: URL query params (flags)

	// dropProblems()

	// fetchAndParseProblems()

	// getAllProblems()

	testFilters()

}

func dropProblems() {
	database, err := db.CreateDB()
	must(err)
	database.DropAllProblems()
}

func testFilters() {
	database, err := db.CreateDB()
	must(err)
	probs, err := database.GetAllProblems()

	// err = database.SetProblemCompleted(689)
	// must(err)
	err = database.SetQuestionBad(689)
	must(err)

	easyProbs := problem.FilterByDifficulty(probs, problem.HARD)
	filtered := problem.FilterByTopic(easyProbs, "array")
	filtered = problem.FilterOutPaid(filtered)
	for _, prob := range filtered {
		fmt.Printf("%v\t%v\t%v\n", prob.DisplayId, prob.Name, prob.BadQuestion)
	}
	fmt.Printf("%d filtered problems.\n", len(filtered))
}

func getAllProblems() {
	database, err := db.CreateDB()
	must(err)
	probs, err := database.GetAllProblems()
	must(err)
	for _, prob := range probs {
		fmt.Printf("%v\t%v\n", prob.DisplayId, prob.Name)
	}
	fmt.Printf("Fetched %d problems.\n", len(probs))
}

func fetchAndParseProblems() {
	// Open a connection to Leetcode with the user-specified query params
	httpResp, err := http.Get(leetcodeApiUrl)
	must(err)

	// Get the JSON body from the response
	htmlReader := httpResp.Body
	defer htmlReader.Close()

	// Parse the questions from the JSON
	problems, err := parser.ParseProblems(htmlReader)
	must(err)

	database, err := db.CreateDB()
	must(err)
	for _, problem := range problems {
		database.InsertProblem(problem)
		// fmt.Println(problem.DisplayId)
		// fmt.Printf("%+v\n", problem)
	}

	fmt.Printf("Fetched %d problems.\n", len(problems))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
