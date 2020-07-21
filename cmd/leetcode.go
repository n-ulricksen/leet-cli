package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("vim-go")

	leetcodeUrl := "https://leetcode.com/problemset/all"
	// TODO: URL query params (flags)

	// Open a connection to Leetcode with the user-specified query params
	httpResp, err := http.Get(leetcodeUrl)
	if err != nil {
		panic(err)
	}
	htmlReader := httpResp.Body
	defer htmlReader.Close()

	questions, err := parser.ParseQuestions(htmlReader)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", questions)
}
