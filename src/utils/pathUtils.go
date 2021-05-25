package utils

import (
	"path/filepath"

	"github.com/skmonir/mango/src/config"
)

func GetSourceDirPath(cfg config.Configuration) string {
	return filepath.Join(cfg.Workspace, cfg.OJ, cfg.CurrentContestId, "src")
}

func GetSourceFilePathWithExt(cfg config.Configuration, problemId string) string {
	return filepath.Join(GetSourceDirPath(cfg), problemId+".cpp")
}

func GetSourceFilePathWithoutExt(cfg config.Configuration, problemId string) string {
	return filepath.Join(GetSourceDirPath(cfg), problemId)
}

func GetTestcaseDirPath(cfg config.Configuration) string {
	return filepath.Join(cfg.Workspace, cfg.OJ, cfg.CurrentContestId, "testcase")
}

func GetTestcaseFilePath(cfg config.Configuration, problemId string) string {
	return filepath.Join(GetTestcaseDirPath(cfg), problemId+".json")
}

func ResolveTescasePath(cfg config.Configuration, problemId string) error {
	testCaseDirPath := GetTestcaseDirPath(cfg)

	if err := CreateFile(testCaseDirPath, problemId+".json"); err != nil {
		return err
	}

	return nil
}
