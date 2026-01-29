package problem

type Problem struct {
	ID         int        `json:"id"`
	Question   string     `json:"question"`
	Difficulty Difficulty `json:"difficulty"`
}

type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)