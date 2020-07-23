package main

import (
	"fmt"
	"net/http"

	"github.com/ulricksennick/leetcode-fetcher/parser"
)

const (
	leetcodeApiUrl = "https://leetcode.com/api/problems/all/"
)

func main() {
	fmt.Println("vim-go")

	// TODO: URL query params (flags)

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

	fmt.Printf("Fetched %d problems.\n", len(problems))

}
