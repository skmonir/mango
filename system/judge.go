package system

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/k0kubun/go-ansi"
	"github.com/shirou/gopsutil/process"
	"github.com/skmonir/mango/config"
	"github.com/skmonir/mango/dto"
	"github.com/skmonir/mango/utils"
	// "syscall"
	// "io"
)

func getCompilationCommand(cfg config.Configuration, problemId string) (string, error) {
	filePathWithExt := utils.GetSourceFilePathWithExt(cfg, problemId)
	filePathWithoutExt := utils.GetSourceFilePathWithoutExt(cfg, problemId)

	if !utils.IsFileExist(filePathWithExt) {
		return "", errors.New("source file not found")
	}

	command := fmt.Sprintf("%v %v %v -o %v", cfg.CompilationCommand, cfg.CompilationArgs, filePathWithExt, filePathWithoutExt)

	return command, nil
}

func getExecutionCommand(cfg config.Configuration, problemId string) string {
	command := utils.GetSourceFilePathWithoutExt(cfg, problemId)
	return command
}

func Compile(cfg config.Configuration, problemId string) error {
	command, err := getCompilationCommand(cfg, problemId)
	if err != nil {
		return err
	}

	cmds := utils.ParseCommand(command)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, cmds[0], cmds[1:]...)
	// cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return errors.New("compile error")
	}

	return nil
}

func Execute(testcase dto.Testcase, command string) dto.ExecutionResult {
	var response dto.ExecutionResult
	response.Test = testcase

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(testcase.TimeLimit+1)*time.Second)
	defer cancel()

	input_buffer := bytes.NewBuffer([]byte(testcase.Input))
	var output_buffer bytes.Buffer

	cmd := exec.CommandContext(ctx, command)
	// cmd.Stderr = os.Stderr
	cmd.Stdin = input_buffer
	cmd.Stdout = &output_buffer

	maxMemory := uint64(0)

	completeExecution := func(err error) {
		response.Status = (err == nil)
		response.Error = err
		response.Memory = maxMemory
		response.Runtime = cmd.ProcessState.UserTime().Milliseconds()
	}

	// if err := cmd.Run(); err != nil {
	// 	response.Status = false
	// 	response.Error = err
	// } else {
	// 	response.Status = true
	// }

	// response.ExitCode = cmd.ProcessState.ExitCode()
	// response.Memory = uint64(cmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss)
	// response.Runtime = cmd.ProcessState.UserTime().Milliseconds()

	if err := cmd.Start(); err != nil {
		completeExecution(err)
		return response
	}

	pid := int32(cmd.Process.Pid)
	ch := make(chan error)
	go func() {
		ch <- cmd.Wait()
	}()
	running := true
	for running {
		select {
		case err := <-ch:
			completeExecution(err)
			if err != nil {
				return response
			}
			running = false
		default:
			p, err := process.NewProcess(pid)
			if err == nil {
				m, err := p.MemoryInfo()
				if err == nil && m.RSS > maxMemory {
					maxMemory = m.RSS
				}
			}
		}
	}

	if !response.Status {
		return response
	}

	response.Output = utils.TrimIO(output_buffer.String())

	return response
}

func Judge(testcase dto.Testcase, executionResult *dto.ExecutionResult) {
	if !executionResult.Status {
		errorMsg := executionResult.Error.Error()
		if strings.Contains(errorMsg, "segmentation fault") {
			executionResult.Verdict = "RE"
		} else if strings.Contains(errorMsg, "killed") {
			executionResult.Verdict = "TLE"
			executionResult.Runtime = testcase.TimeLimit * 1000
		} else if utils.ConvertMemoryInMb(executionResult.Memory) > testcase.MemoryLimit {
			executionResult.Verdict = "MLE"
		} else {
			executionResult.Verdict = "RE"
		}
	} else {
		if executionResult.Runtime > (testcase.TimeLimit * 1000) {
			executionResult.Verdict = "TLE"
			executionResult.Runtime = testcase.TimeLimit * 1000
		} else if utils.ConvertMemoryInMb(executionResult.Memory) > testcase.MemoryLimit {
			executionResult.Verdict = "MLE"
		} else if testcase.Output == executionResult.Output {
			executionResult.Verdict = "OK"
		} else {
			executionResult.Verdict = "WA"
		}
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

func RunTest(cfg config.Configuration, problemId string) error {
	problemInfo, err := utils.GetProblemInfo(cfg, problemId)
	if err != nil {
		ansi.Println(color.RedString("could not fetch problem info"))
		return nil
	}

	if problemInfo.Name != "" {
		ansi.Println(color.CyanString("\nProblem Name: " + problemInfo.Name))
		ansi.Println(color.New(color.FgCyan).Sprintf("Time Limit: %v sec, Memory Limit: %v MB\n", problemInfo.TimeLimit, problemInfo.MemoryLimit))
	}

	// ansi.Printf("Compilation Status: ")
	// ansi.Printf(color.MagentaString("compiling..."))
	if err := Compile(cfg, problemId); err != nil {
		// PublishCompilationResult(err)
		ansi.Println(color.RedString(err.Error()))
		return nil
	}

	testcases := problemInfo.Dataset
	testResults := make([]dto.ExecutionResult, len(testcases))
	command := getExecutionCommand(cfg, problemId)

	for idx, testcase := range testcases {
		executionResult := Execute(testcase, command)
		Judge(testcase, &executionResult)
		testResults[idx] = executionResult
	}

	PublishTestResult(problemInfo, testResults)

	return nil
}
