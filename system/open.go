package system

import (
	"fmt"
	"errors"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/skmonir/mango/config"
	"github.com/skmonir/mango/utils"
)

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func OpenContest(cfg config.Configuration, problemIdList []string) error {
	var err error
	for _, problemId := range problemIdList {
		sourcePath := utils.GetSourceFilePathWithExt(cfg, problemId)

		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", sourcePath).Run()
		case "windows":
			exec.Command("cmd", fmt.Sprintf("/C start %v", sourcePath)).Run()
		case "darwin":
			err = exec.Command("open", sourcePath).Run()
		default:
			ansi.Println(color.New(color.FgRed).Sprintf("unsupported os"))
		}
	}
	return err
}

func OpenProblem(cfg config.Configuration, problemId string) error {
	if err := OpenContest(cfg, []string{problemId}); err != nil {
		return errors.New("error while opening source")
	}
	return nil
}

func Open(cfg config.Configuration, cmd string) error {
	contestId, problemId, err := utils.ParseContestAndProblemId(cmd)
	if err != nil {
		return err
	}
	if contestId != "" {
		cfg.CurrentContestId = contestId
	}
	if cfg.CurrentContestId == "" {
		return errors.New("please set current contest id or use contest & problem id combination like 1512G")
	}

	if problemId == "" {
		sourcePath := utils.GetSourceDirPath(cfg)
		problemIdList := utils.GetFilenamesInDir(sourcePath)
		if err := OpenContest(cfg, problemIdList); err != nil {
			return err
		}
	} else {
		if err := OpenProblem(cfg, problemId); err != nil {
			return err
		}
	}

	return nil
}
