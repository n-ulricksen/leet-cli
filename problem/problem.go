package problem

const (
	EASY   int = 1
	MEDIUM int = 2
	HARD   int = 3
)

type Problem struct {
	Name       string   `json:"name"`
	Id         int      `json:"id"`
	Url        string   `json:"url"`
	Difficulty int      `json:"difficulty"`
	Topics     []string `json:"topics"`
	Paid       bool     `json:"paid"`
	Upvotes    int      `json:"upvote"`
	Downvotes  int      `json:"downvotes"`
	Acceptance float32  `json:"acceptance"`
	Completed  bool     `json:"completed"`
}
