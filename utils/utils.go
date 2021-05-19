package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"unicode"

	"github.com/skmonir/mango/config"
	"github.com/skmonir/mango/dto"
)

func IsDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func IsFileExist(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func IsDirExist(folderPath string) bool {
	_, err := os.Stat(folderPath)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDir(folderPath string) error {
	if !IsDirExist(folderPath) {
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func GetFilenamesInDir(folderPath string) []string {
	var filenames []string

	if !IsDirExist(folderPath) {
		return filenames
	}

	files, err := ioutil.ReadDir(folderPath)

	if err == nil {
		for _, f := range files {
			fname := f.Name()
			if strings.HasSuffix(fname, ".cpp") {
				fname = strings.TrimSuffix(fname, filepath.Ext(fname))
				filenames = append(filenames, fname)
			}
		}
	}
	return filenames
}

func CreateFile(folderPath string, filename string) error {
	filePath := filepath.Join(folderPath, filename)
	if !IsFileExist(filePath) {
		if err := CreateDir(folderPath); err != nil {
			return err
		}
		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer f.Close()
	}
	return nil
}

func OpenFile(filePath string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", filePath).Run()
	case "windows":
		exec.Command("cmd", filePath).Run()
	case "darwin":
		err = exec.Command("open", filePath).Run()
	default:
		err = errors.New("unsupported os")
	}

	return err
}

// ParseCommand parses a command line and handle arguments in quotes.
// https://github.com/vrischmann/shlex/blob/master/shlex.go
func ParseCommand(s string) (res []string) {
	var buf bytes.Buffer
	insideQuotes := false
	for _, r := range s {
		switch {
		case unicode.IsSpace(r) && !insideQuotes:
			if buf.Len() > 0 {
				res = append(res, buf.String())
				buf.Reset()
			}
		case r == '"' || r == '\'':
			if insideQuotes {
				res = append(res, buf.String())
				buf.Reset()
				insideQuotes = false
				continue
			}
			insideQuotes = true
		default:
			buf.WriteRune(r)
		}
	}
	if buf.Len() > 0 {
		res = append(res, buf.String())
	}
	return
}

func ParseContestAndProblemId(cmd string) (string, string, error) {
	cmd = strings.TrimSpace(cmd)
	if len(cmd) == 0 {
		return "", "", errors.New("command is not valid")
	}

	ptr, sz := 0, len(cmd)

	contestId := ""
	for ptr < sz && IsDigit(rune(cmd[ptr])) {
		contestId += string(cmd[ptr])
		ptr++
	}

	problemId := ""
	for ptr < sz && rune(cmd[ptr]) != ' ' {
		problemId += string(cmd[ptr])
		ptr++
	}

	return contestId, problemId, nil
}

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

func GetProblemInfo(cfg config.Configuration, problemId string) (dto.Problem, error) {
	var problemInfo dto.Problem
	testpath := GetTestcaseFilePath(cfg, problemId)

	data, err := ioutil.ReadFile(testpath)
	if err != nil {
		return problemInfo, err
	}

	err = json.Unmarshal(data, &problemInfo)
	if err != nil {
		return problemInfo, err
	}

	return problemInfo, nil
}

func ConvertMemoryInMb(memory uint64) uint64 {
	return memory / 1024 / 1024
}

func ParseMemoryInMb(memory uint64) string {
	return fmt.Sprintf("%v MB", memory/1024/1024)
}

func ParseMemoryInKb(memory uint64) string {
	return fmt.Sprintf("%v KB", memory/1024)
}

func GetContestUrl(cfg config.Configuration) string {
	return fmt.Sprintf("%v/contest/%v", cfg.Host, cfg.CurrentContestId)
}

func GetProblemUrl(cfg config.Configuration, problemId string) string {
	return fmt.Sprintf("%v/contest/%v/problem/%v", cfg.Host, cfg.CurrentContestId, problemId)
}

func FilterHtml(src []byte) []byte {
	newline := regexp.MustCompile(`<[\s/br]+?>`)
	src = newline.ReplaceAll(src, []byte("\n"))
	s := html.UnescapeString(string(src))
	return []byte(s)
}

func TrimIO(io string) string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(io))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	io = strings.Join(lines, "\n")
	return io
}
