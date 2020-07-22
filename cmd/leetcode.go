package main

import (
	"fmt"
	"net/http"

	"github.com/ulricksennick/leetcode-fetcher"
)

func main() {
	fmt.Println("vim-go")

	leetcodeUrl := "https://leetcode.com/api/problems/all/"
	// TODO: URL query params (flags)

	// Open a connection to Leetcode with the user-specified query params
	httpResp, err := http.Get(leetcodeUrl)
	if err != nil {
		panic(err)
	}

	// Get the JSON body from the response
	htmlReader := httpResp.Body
	defer htmlReader.Close()

	// Parse the questions from the JSON
	questions, err := parser.ParseQuestions(htmlReader)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", questions[0])
	fmt.Printf("Fetched %d questions.\n", len(questions))
}
