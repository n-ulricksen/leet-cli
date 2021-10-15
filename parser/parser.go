package parser

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ulricksennick/leet-cli/db"
	"github.com/ulricksennick/leet-cli/problem"
	"github.com/ulricksennick/leet-cli/urls"
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

type tmpTopics struct {
	Topics []struct {
		Slug      string `json:"slug"`
		Questions []int  `json:"questions"`
	} `json:"topics"`
}

func FetchAndStoreProblems() {
	htmlReader := getHttpBody(urls.LeetcodeApiUrl)
	defer htmlReader.Close()

	problems, err := parseProblems(htmlReader)
	must(err)

	database, err := db.CreateDB()
	err = database.InsertProblems(problems)
	must(err)
}

func FetchAndStoreTopics() {
	topics := fetchTopics()

	database, err := db.CreateDB()
	database.DropAllTopics()
	err = database.InsertTopics(topics)
	must(err)
}

func parseProblems(r io.Reader) (map[int]*problem.Problem, error) {
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
			Slug:       q.Stat.Slug,
			Url:        urls.LeetcodeProblemUrl + q.Stat.Slug,
			Difficulty: q.Difficulty["level"],
			Paid:       q.Paid,
		}
	}

	updateProblemTopics(problems)

	return problems, nil
}

func fetchTopics() []*problem.Topic {
	httpBody := getHttpBody(urls.ProblemTopicUrl)
	defer httpBody.Close()

	slugsTmp := new(tmpTopics)
	json.NewDecoder(httpBody).Decode(slugsTmp)

	topics := make([]*problem.Topic, len(slugsTmp.Topics))
	for i, t := range slugsTmp.Topics {
		newTopic := &problem.Topic{
			Slug:      t.Slug,
			Questions: t.Questions,
		}
		topics[i] = newTopic
	}

	return topics
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// Get the list of problems/topics, assign topics to appropriate problems
func updateProblemTopics(problems map[int]*problem.Problem) {
	if len(problems) == 0 {
		return
	}

	database, err := db.CreateDB()
	must(err)

	topics, err := database.GetAllTopics()
	must(err)

	for _, topic := range topics {
		for _, questionId := range topic.Questions {
			if problems[questionId] == nil {
				continue
			}
			problems[questionId].Topics = append(problems[questionId].Topics, topic.Slug)
		}
	}
}

func getHttpBody(url string) io.ReadCloser {
	httpResp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	return httpResp.Body
}
