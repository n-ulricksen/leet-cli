package parser

import (
	"io"

	"golang.org/x/net/html"
)

type Question struct {
	Name       string
	URL        string
	Difficulty string
	Upvotes    int
	Downvotes  int
	Acceptence float32
}

const (
	EASY   = "Easy"
	MEDIUM = "Medium"
	HARD   = "Hard"
)

func ParseQuestions(r io.Reader) ([]Question, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var questions []Question

	var dfs func(*html.Node)
	dfs = func(n *html.Node) {
		// Find for leetcode question

		for c := n.FirstChild; c != nil; c = n.NextSibling {
			dfs(c)
		}
	}
	dfs(doc)

	return questions, nil
}
