package parser

import (
	"encoding/json"
	"io"
)

type tmpQuestion struct {
	Stat struct {
		Name string `json:"question__title"`
		Id   int    `json:"question_id"`
		Slug string `json:"question__title_slug"`
	} `json:"stat"`
	Difficulty map[string]int `json:"difficulty"`
	Paid       bool           `json:"paid_only"`
}
type tmpQuestions struct {
	Questions []tmpQuestion `json:"stat_status_pairs"`
}

type Question struct {
	Name       string
	Id         int
	Slug       string
	Difficulty int
	Paid       bool
	Upvotes    int
	Downvotes  int
	Acceptance float32
}

func ParseQuestions(r io.Reader) ([]Question, error) {
	// Decode the JSON
	var tmp tmpQuestions
	err := json.NewDecoder(r).Decode(&tmp)
	if err != nil {
		panic(err)
	}

	// TODO: navigate to question's URL to find upvotes, downvotes, acceptance

	// Create well formatted questions
	var questions []Question
	for _, q := range tmp.Questions {
		questions = append(questions, Question{
			Name:       q.Stat.Name,
			Id:         q.Stat.Id,
			Slug:       q.Stat.Slug,
			Difficulty: q.Difficulty["level"],
			Paid:       q.Paid,
		})
	}

	return questions, nil
}
