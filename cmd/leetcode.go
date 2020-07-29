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

	// database.DropAllProblems()

	// fetchAndParseProblems()

	// getAllProblems()
	// getEasyProblems()
	// getTopicProblems("two-pointers")

	database, err := db.CreateDB()
	must(err)
	probs, err := database.GetAllProblems()

	easyProbs := problem.FilterByDifficulty(probs, problem.HARD)
	filtered := problem.FilterByTopic(easyProbs, "array")
	filtered = problem.FilterOutPaid(filtered)
	for _, e := range filtered {
		fmt.Printf("%+v\n", e)
	}
	fmt.Println(len(filtered))
}

func getAllProblems() {
	database, err := db.CreateDB()
	must(err)
	probs, err := database.GetAllProblems()
	must(err)
	for _, prob := range probs {
		fmt.Printf("%+v\n", *prob)
	}
	fmt.Printf("Fetched %d problems.\n", len(probs))
}

func getEasyProblems() {
	database, err := db.CreateDB()
	must(err)
	probs, err := database.GetProblemsByDifficulty(problem.EASY)
	must(err)
	for _, prob := range probs {
		fmt.Printf("%+v\n", *prob)
	}
	fmt.Printf("Fetched %d problems.\n", len(probs))
}

func getTopicProblems(topic string) {
	database, err := db.CreateDB()
	must(err)
	probs, err := database.GetProblemsByTopic(topic)
	must(err)
	for _, prob := range probs {
		fmt.Printf("%+v\n", *prob)
	}
	fmt.Printf("Fetched %d problems.\n", len(probs))
}

func fetchAndParseProblems() {
	// Open a connection to Leetcode with the user-specified query params
	httpResp, err := http.Get(leetcodeApiUrl)
	if err != nil {
		panic(err)
	}

	// Get the JSON body from the response
	htmlReader := httpResp.Body
	defer htmlReader.Close()

	// Parse the questions from the JSON
	problems, err := parser.ParseProblems(htmlReader)
	if err != nil {
		panic(err)
	}

	database, err := db.CreateDB()
	must(err)
	for _, problem := range problems {
		database.InsertProblem(problem)
		fmt.Printf("%+v\n", problem)
	}

	fmt.Printf("Fetched %d problems.\n", len(problems))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
