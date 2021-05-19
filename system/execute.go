package system

import (
	"bytes"
	"context"
	"os/exec"
	"time"

	"github.com/shirou/gopsutil/process"
	"github.com/skmonir/mango/config"
	"github.com/skmonir/mango/dto"
	"github.com/skmonir/mango/utils"
)

func getExecutionCommand(cfg config.Configuration, problemId string) string {
	command := utils.GetSourceFilePathWithoutExt(cfg, problemId)
	return command
}

func ExecuteSourceBinary(cfg config.Configuration, testcase dto.Testcase, problemId string) dto.ExecutionResult {
	command := getExecutionCommand(cfg, problemId)

	var response dto.ExecutionResult
	response.Test = testcase

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(testcase.TimeLimit)*time.Second)
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
