package problem

const (
	EASY   int = 1
	MEDIUM int = 2
	HARD   int = 3
)

var DIFFICULTY_STRINGS = []string{"easy", "medium", "hard"}

type Problem struct {
	Name           string   `json:"name"`
	Id             int      `json:"id"`
	DisplayId      int      `json:"displayId"` // identify a problem on Leetcode website
	Url            string   `json:"url"`
	Slug           string   `json:"titleSlug"`
	Difficulty     int      `json:"difficulty"`
	Topics         []string `json:"topics"`
	Paid           bool     `json:"paid"`
	Completed      bool     `json:"completed"`
	IsBad          bool     `json:"isBad"`
	SampleTestCase string   `json:"sampleTestCase"`
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

// Return problem set containing only incomplete problems.
func FilterOutCompleted(problems []*Problem) []*Problem {
	var ret []*Problem

	for _, prob := range problems {
		if !prob.Completed {
			ret = append(ret, prob)
		}
	}

	return ret
}

// Return problem set containing only completed problems.
func FilterCompleted(problems []*Problem) []*Problem {
	var ret []*Problem

	for _, prob := range problems {
		if prob.Completed {
			ret = append(ret, prob)
		}
	}

	return ret
}
