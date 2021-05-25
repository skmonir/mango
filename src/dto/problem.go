package dto

type Problem struct {
	Name        string
	TimeLimit   int64
	MemoryLimit uint64
	Dataset     []Testcase
}
