package system

import (
	"errors"

	"github.com/skmonir/mango/config"
	"github.com/skmonir/mango/utils"
)

func CreateProblem(cfg config.Configuration, problemId string) error {
	if err := ParseProblem(cfg, problemId); err != nil {
		return err
	}

	if err := CreateSource(cfg, problemId); err != nil {
		return err
	}
	CopyTemplateToSource(cfg, problemId)
	OpenProblem(cfg, problemId)
	return nil
}

func CreateContest(cfg config.Configuration) error {
	URL := utils.GetContestUrl(cfg)
	problemIdList, err := GetProblemIdList(URL)
	if err != nil {
		return err
	}

	if err = ParseContest(cfg, problemIdList); err != nil {
		return err
	}

	CreateSourceList(cfg, problemIdList)
	CopyTemplateToSourceList(cfg, problemIdList)
	OpenContest(cfg, problemIdList)

	return err
}

func Create(cfg config.Configuration, cmd string) error {
	contestId, problemId, err := utils.ParseContestAndProblemId(cmd)
	if err != nil {
		return err
	}

	cfg.CurrentContestId = contestId
	if cfg.CurrentContestId == "" {
		return errors.New("please use contest & problem id combination like 1512G")
	}

	if problemId == "" {
		if err := CreateContest(cfg); err != nil {
			return err
		}
	} else {
		if err := CreateProblem(cfg, problemId); err != nil {
			return err
		}
	}

	if err := config.SetContest(contestId); err != nil {
		return errors.New("error while saving config")
	}

	return nil
}
