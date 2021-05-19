package system

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/skmonir/mango/config"
	"github.com/skmonir/mango/utils"
)

func CopyTemplateToSource(cfg config.Configuration, problemId string) error {
	if cfg.TemplatePath != "" {
		srcPath := cfg.TemplatePath
		destPath := utils.GetSourceFilePathWithExt(cfg, problemId)
		problemInfo, _ := utils.GetProblemInfo(cfg, problemId)

		if !utils.IsFileExist(srcPath) {
			return errors.New("template file not found")
		}
		if !utils.IsFileExist(destPath) {
			return errors.New("source file not found")
		}

		template := ""

		srcFile, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		scanner := bufio.NewScanner(srcFile)
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimRight(line, " ")
			line = strings.ReplaceAll(line, "%author%", cfg.Author)
			line = strings.ReplaceAll(line, "%problem%", problemInfo.Name)
			line = strings.ReplaceAll(line, "%datetime%", time.Now().Local().Format("2-Jan-2006 15:04:05"))
			template += line + "\n"
		}

		if err = scanner.Err(); err != nil {
			return err
		}

		if err := ioutil.WriteFile(destPath, []byte(template), 0644); err != nil {
			return err
		}
	}

	return nil
}

func CopyTemplateToSourceList(cfg config.Configuration, problemIdList []string) error {
	var err error = nil
	for _, problemId := range problemIdList {
		err = CopyTemplateToSource(cfg, problemId)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return err
}

func CreateSource(cfg config.Configuration, problemId string) error {
	sourcePath := utils.GetSourceDirPath(cfg)
	sourceName := problemId + ".cpp"
	if err := utils.CreateFile(sourcePath, sourceName); err != nil {
		return err
		// ansi.Println(color.New(color.FgRed).Sprintf("error while creating source for task %v", problemId))
	}
	return nil
}

func CreateSourceList(cfg config.Configuration, problemIdList []string) error {
	var err error = nil
	for _, problemId := range problemIdList {
		err = CreateSource(cfg, problemId)
	}
	return err
}

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
