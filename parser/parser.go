package parser

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ulricksennick/lcfetch/problem"
)

const (
	leetcodeBaseUrl = "https://leetcode.com/problems/"
	problemTopicUrl = "https://leetcode.com/problems/api/tags/"
	leetcodeApiUrl  = "https://leetcode.com/api/problems/all/"
)

type tmpProblems struct {
	Problems []struct {
		Stat struct {
			Name      string `json:"question__title"`
			Id        int    `json:"question_id"`
			DisplayId int    `json:"frontend_question_id"`
			Slug      string `json:"question__title_slug"`
		} `json:"stat"`
		Difficulty map[string]int `json:"difficulty"`
		Paid       bool           `json:"paid_only"`
	} `json:"stat_status_pairs"`
}

type slugs struct {
	Topics []struct {
		Slug      string `json:"slug"`
		Questions []int  `json:"questions"`
	} `json:"topics"`
}

func ParseProblems(r io.Reader) (map[int]*problem.Problem, error) {
	// Decode the JSON
	decoded := new(tmpProblems)
	err := json.NewDecoder(r).Decode(decoded)
	if err != nil {
		return nil, err
	}

	// Create well formatted problems
	problems := make(map[int]*problem.Problem)
	for _, q := range decoded.Problems {
		id := q.Stat.Id

		problems[id] = &problem.Problem{
			Name:       q.Stat.Name,
			Id:         q.Stat.Id,
			DisplayId:  q.Stat.DisplayId,
			Url:        leetcodeBaseUrl + q.Stat.Slug,
			Difficulty: q.Difficulty["level"],
			Paid:       q.Paid,
		}
	}

	updateProblemTopics(problems)

	return problems, nil
}

// Get the list of problems/topics, assign topics to appropriate problems
func updateProblemTopics(problems map[int]*problem.Problem) {
	if len(problems) == 0 {
		return
	}

	httpBody := getHttpBody(problemTopicUrl)
	defer httpBody.Close()

	slugsTmp := new(slugs)
	json.NewDecoder(httpBody).Decode(slugsTmp)

	for _, topic := range slugsTmp.Topics {
		for _, questionId := range topic.Questions {
			if problems[questionId] == nil {
				continue
			}
			problems[questionId].Topics = append(problems[questionId].Topics, topic.Slug)
		}
	}
}

func getHttpBody(url string) io.ReadCloser {
	httpResp, err := http.Get(problemTopicUrl)
	if err != nil {
		panic(err)
	}
	return httpResp.Body
}
