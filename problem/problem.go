package problem

type Problem struct {
	Name       string
	Id         int
	Url        string
	Difficulty int
	Topics     []string
	Paid       bool
	Upvotes    int
	Downvotes  int
	Acceptance float32
}
