package dto

type ExecutionResult struct {
	Status   bool
	Message  string
	Error    error
	Output   string
	Verdict  string
	Runtime  int64
	Memory   uint64
	ExitCode int
	Test     Testcase
}
