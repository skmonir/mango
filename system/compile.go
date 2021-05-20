package system

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/skmonir/mango/config"
	"github.com/skmonir/mango/utils"
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

func CompileSource(cfg config.Configuration, problemId string, showStdError bool) error {
	command, err := getCompilationCommand(cfg, problemId)
	if err != nil {
		return err
	}

	cmds := utils.ParseCommand(command)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, cmds[0], cmds[1:]...)
	if showStdError {
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return errors.New("compile error")
	}

	return nil
}

func Compile(cfg config.Configuration, cmd string) error {
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
		return errors.New("please set contest & problem id combination like 1512G")
	}

	if err := CompileSource(cfg, problemId, true); err != nil {
		return err
	}

	return nil
}
