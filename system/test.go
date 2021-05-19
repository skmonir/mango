package system

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/k0kubun/go-ansi"
	"github.com/skmonir/mango/config"
	"github.com/skmonir/mango/dto"
	"github.com/skmonir/mango/utils"
)

func GetVerdict(testcase dto.Testcase, executionResult *dto.ExecutionResult) {
	errorMsg := ""
	if !executionResult.Status && executionResult.Error != nil {
		errorMsg = executionResult.Error.Error()
	}

	if strings.Contains(errorMsg, "segmentation fault") {
		executionResult.Verdict = "RE"
	} else if (executionResult.Runtime > (testcase.TimeLimit * 1000)) || strings.Contains(errorMsg, "killed") {
		executionResult.Verdict = "TLE"
		executionResult.Runtime = testcase.TimeLimit * 1000
	} else if utils.ConvertMemoryInMb(executionResult.Memory) > testcase.MemoryLimit {
		executionResult.Verdict = "MLE"
	} else if !executionResult.Status {
		executionResult.Verdict = "RE"
	} else if testcase.Output == executionResult.Output {
		executionResult.Verdict = "OK"
	} else {
		executionResult.Verdict = "WA"
	}
}

func PublishTestResult(problemInfo dto.Problem, executionResultList []dto.ExecutionResult) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"INPUT", "OUPUT", "EXPECTED", "VERDICT", "TIME", "MEMORY"})

	totalPassed, totalTests := 0, len(executionResultList)
	for idx, test := range executionResultList {
		if idx > 0 {
			t.AppendSeparator()
		}
		runtime := fmt.Sprintf("%v ms", test.Runtime)
		memory := utils.ParseMemoryInKb(test.Memory)
		verdict := test.Verdict + "\n"

		if test.Verdict == "OK" {
			totalPassed++
		}

		t.AppendRow([]interface{}{test.Test.Input, test.Output, test.Test.Output, verdict, runtime, memory})
	}
	t.SetStyle(table.StyleLight)
	t.Render()

	if totalTests != totalPassed {
		ansi.Println(color.New(color.FgRed).Sprintf("\n%v/%v Tests Passed\n", totalPassed, totalTests))
	} else {
		ansi.Println(color.New(color.FgGreen).Sprintf("\n%v/%v Tests Passed\n", totalPassed, totalTests))
	}
}

func RunTest(cfg config.Configuration, cmd string) error {
	contestId, problemId, err := utils.ParseContestAndProblemId(cmd)
	if err != nil {
		return err
	}
	if problemId == "" {
		return errors.New("problem id not valid")
	}
	if contestId != "" {
		cfg.CurrentContestId = contestId
	}
	if cfg.CurrentContestId == "" {
		return errors.New("please set current contest id or use contest & problem id combination like 1512G")
	}

	problemInfo, err := utils.GetProblemInfo(cfg, problemId)
	if err != nil {
		ansi.Println(color.RedString("could not fetch problem info"))
		return nil
	}

	if problemInfo.Name != "" {
		ansi.Println(color.CyanString("\nProblem Name: " + problemInfo.Name))
		ansi.Println(color.New(color.FgCyan).Sprintf("Time Limit: %v sec, Memory Limit: %v MB\n", problemInfo.TimeLimit, problemInfo.MemoryLimit))
	}

	if err := CompileSource(cfg, problemId); err != nil {
		ansi.Println(color.RedString(err.Error()))
		return nil
	}

	testcases := problemInfo.Dataset
	testResults := make([]dto.ExecutionResult, len(testcases))

	for idx, testcase := range testcases {
		executionResult := ExecuteSourceBinary(cfg, testcase, problemId)
		GetVerdict(testcase, &executionResult)
		testResults[idx] = executionResult
	}

	PublishTestResult(problemInfo, testResults)

	return nil
}
