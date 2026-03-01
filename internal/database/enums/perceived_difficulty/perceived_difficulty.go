package perceived_difficulty

type PerceivedDifficulty int

const (
	VeryEasy PerceivedDifficulty = iota
	Easy
	Medium
	Hard
	VeryHard
	ExtremelyHard
)