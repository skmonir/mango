package utils

import (
	"bytes"
	"errors"
	"strings"
	"unicode"
)

func IsDigit(r rune) bool {
	return '0' <= r && r <= '9'
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
