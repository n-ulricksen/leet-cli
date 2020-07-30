package problem

const (
	EASY   int = 1
	MEDIUM int = 2
	HARD   int = 3
)

type Problem struct {
	Name       string   `json:"name"`
	Id         int      `json:"id"`
	DisplayId  int      `json:"displayId"` // identify problem on Leetcode website
	Url        string   `json:"url"`
	Difficulty int      `json:"difficulty"`
	Topics     []string `json:"topics"`
	Paid       bool     `json:"paid"`
	Upvotes    int      `json:"upvote"`
	Downvotes  int      `json:"downvotes"`
	Acceptance float32  `json:"acceptance"`
	Completed  bool     `json:"completed"`
	IsBad      bool     `json:"isBad"`
}

func FilterByDifficulty(problems []*Problem, difficulty int) []*Problem {
	var ret []*Problem

	for _, prob := range problems {
		if prob.Difficulty == difficulty {
			ret = append(ret, prob)
		}
	}

	return ret
}

func FilterByTopic(problems []*Problem, topic string) []*Problem {
	var ret []*Problem

	for _, prob := range problems {
		for _, t := range prob.Topics {
			if t == topic {
				ret = append(ret, prob)
			}
		}
	}

	return ret
}

func FilterOutPaid(problems []*Problem) []*Problem {
	var ret []*Problem

	for _, prob := range problems {
		if !prob.Paid {
			ret = append(ret, prob)
		}
	}

	return ret
}

func FilterOutComplted(problems []*Problem) []*Problem {
	var ret []*Problem

	for _, prob := range problems {
		if !prob.Completed {
			ret = append(ret, prob)
		}
	}

	return ret
}
