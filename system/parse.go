package system

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/skmonir/mango/config"
	"github.com/skmonir/mango/dto"
	"github.com/skmonir/mango/utils"
)

// Implementation idea's are collected from https://github.com/xalanq/cf-tool/blob/master/client/parse.go
// which is a similar project and has more features. Methods are customized as per need.

func GetBody(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func GetProblemIdList(URL string) ([]string, error) {
	body, err := GetBody(URL)
	if err != nil {
		return []string{}, errors.New("error while problem list")
	}

	stat_regexp := regexp.MustCompile(`class="problems"[\s\S]+?</tr>([\s\S]+?)</table>`)
	stat_body := stat_regexp.FindSubmatch(body)
	if stat_body == nil {
		return []string{}, errors.New("no problems found")
	}
	probs_table := stat_body[1]

	row_idx_regex := regexp.MustCompile(`<tr[\s\S]*?>`)
	row_idxs := row_idx_regex.FindAllIndex(probs_table, -1)
	if row_idxs == nil {
		return []string{}, errors.New("no problems found")
	}

	row_idxs = append(row_idxs, []int{0, len(probs_table)})
	id_td_regex := regexp.MustCompile(`<td[\s\S]*?>`)
	prob_id_regex := regexp.MustCompile(`<a[\s\S]*?>([\s\S]*)</a>`)

	problem_ids := make([]string, len(row_idxs)-1)
	for i := 1; i < len(row_idxs); i++ {
		current_row := probs_table[row_idxs[i-1][0]:row_idxs[i][1]]
		td_idxs := id_td_regex.FindAllIndex(current_row, -1)

		current_prob_elem := current_row[td_idxs[0][0]:td_idxs[1][1]]
		id := prob_id_regex.FindSubmatch(current_prob_elem)
		if id != nil {
			problem_ids[i-1] = strings.TrimSpace(string(id[1]))
		} else {
			problem_ids[i-1] = "$"
		}
	}

	return problem_ids, nil
}

func findConstraints(body []byte) (int64, uint64) {
	trg := regexp.MustCompile(`class="time-limit"[\s\S]*?([\d]+) seconds`)
	mrg := regexp.MustCompile(`class="memory-limit"[\s\S]*?([\d]+) megabytes`)
	a := trg.FindSubmatch(body)
	b := mrg.FindSubmatch(body)

	var timeLimit int64 = 2      // default time-limit
	var memoryLimit uint64 = 512 // default memory-limit

	if len(a) > 0 {
		TL, err := strconv.Atoi(strings.TrimSpace(string(utils.FilterHtml(a[1]))))
		if err == nil {
			timeLimit = int64(TL)
		}
	}
	if len(b) > 0 {
		ML, err := strconv.Atoi(strings.TrimSpace(string(utils.FilterHtml(b[1]))))
		if err == nil {
			memoryLimit = uint64(ML)
		}
	}

	return timeLimit, memoryLimit
}

func findSample(body []byte) (input [][]byte, output [][]byte, err error) {
	irg := regexp.MustCompile(`class="input"[\s\S]*?<pre>([\s\S]*?)</pre>`)
	org := regexp.MustCompile(`class="output"[\s\S]*?<pre>([\s\S]*?)</pre>`)
	a := irg.FindAllSubmatch(body, -1)
	b := org.FindAllSubmatch(body, -1)
	if a == nil || b == nil || len(a) != len(b) {
		return nil, nil, fmt.Errorf("cannot parse samples")
	}
	for i := 0; i < len(a); i++ {
		input = append(input, utils.FilterHtml(a[i][1]))
		output = append(output, utils.FilterHtml(b[i][1]))
	}
	return
}

func GetProblemName(body []byte) string {
	name_body_regex := regexp.MustCompile(`class="title"([\s\S]*?)class="time-limit"`)
	name_body := name_body_regex.FindSubmatch(body)
	if name_body == nil {
		return ""
	}

	name_regex := regexp.MustCompile(`>([\s\S]*?)</div>[\s\S]*?`)
	name := name_regex.FindSubmatch(name_body[1])
	if name == nil {
		return ""
	}
	return strings.TrimSpace(string(name[1]))
}

// ParseProblem parse problem to path
func ParseProblem(cfg config.Configuration, problemId string) error {
	URL := utils.GetProblemUrl(cfg, problemId)
	body, err := GetBody(URL)
	if err != nil {
		return err
	}

	probName := GetProblemName(body)
	timeLimit, memoryLimit := findConstraints(body)

	input, output, err := findSample(body)
	if err != nil {
		return err
	}

	// standardIO = true
	// if !bytes.Contains(body, []byte(`<div class="input-file"><div class="property-title">input</div>standard input</div><div class="output-file"><div class="property-title">output</div>standard output</div>`)) {
	// 	standardIO = false
	// }

	testCaseList := make([]dto.Testcase, len(input))
	for i := 0; i < len(input); i++ {
		testCaseList[i] = dto.Testcase{
			Input:       string(input[i]),
			Output:      string(output[i]),
			TimeLimit:   timeLimit,
			MemoryLimit: memoryLimit,
		}
		testCaseList[i].Input = utils.TrimIO(testCaseList[i].Input)
		testCaseList[i].Output = utils.TrimIO(testCaseList[i].Output)
	}

	problem := dto.Problem{
		Name:        probName,
		TimeLimit:   timeLimit,
		MemoryLimit: memoryLimit,
		Dataset:     testCaseList,
	}

	data, err := json.MarshalIndent(problem, "", " ")
	if err != nil {
		return err
	}

	err = utils.ResolveTescasePath(cfg, problemId)
	if err != nil {
		return err
	}

	testCasePath := utils.GetTestcaseFilePath(cfg, problemId)

	err = ioutil.WriteFile(testCasePath, data, 0644)
	if err != nil {
		return err
	}

	ansi.Println(color.New(color.FgGreen).Sprintf("Successfully parsed problem %v", problemId))

	return nil
}

// Parse Contest
func ParseContest(cfg config.Configuration, problemIdList []string) error {
	for _, problemId := range problemIdList {
		if err := ParseProblem(cfg, problemId); err != nil {
			ansi.Println(color.New(color.FgRed).Sprintf("Error while parsing problem %v", problemId))
		}
	}

	return nil
}

func Parse(cfg config.Configuration, cmd string) error {
	contestId, problemId, err := utils.ParseContestAndProblemId(cmd)
	if err != nil {
		return err
	}

	cfg.CurrentContestId = contestId
	if cfg.CurrentContestId == "" {
		return errors.New("please use contest & problem id combination like 1512G")
	}

	if problemId == "" {
		URL := utils.GetContestUrl(cfg)
		problemIdList, err := GetProblemIdList(URL)
		if err != nil {
			return err
		}
		if err := ParseContest(cfg, problemIdList); err != nil {
			return err
		}
	} else {
		if err := ParseProblem(cfg, problemId); err != nil {
			ansi.Println(color.New(color.FgRed).Sprintf("Error while parsing problem %v", problemId))
			return err
		}
	}

	if err := config.SetContest(contestId); err != nil {
		return errors.New("error while saving config")
	}

	return nil
}
