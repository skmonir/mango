package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/skmonir/mango/src/config"
	"github.com/skmonir/mango/src/dto"
)

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

func GetContestType(contestId string) string {
	id, err := strconv.ParseInt(contestId, 10, 30)
	if err != nil {
		if len(contestId) > 5 {
			return "gym"
		}
		return "contest"
	}
	if id > 100000 {
		return "gym"
	}
	return "contest"
}

func GetContestUrl(cfg config.Configuration) string {
	contestType := GetContestType(cfg.CurrentContestId)
	return fmt.Sprintf("%v/%v/%v", cfg.Host, contestType, cfg.CurrentContestId)
}

func GetProblemUrl(cfg config.Configuration, problemId string) string {
	contestType := GetContestType(cfg.CurrentContestId)
	return fmt.Sprintf("%v/%v/%v/problem/%v", cfg.Host, contestType, cfg.CurrentContestId, problemId)
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
