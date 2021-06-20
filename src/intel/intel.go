package intel

import (
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/skmonir/mango/src/config"
)

func WhereAmI() *config.Configuration {
	dir, err := os.Getwd()
	if err != nil {
		return nil
	}
	return whereAmI(dir)
}

func whereAmI(dir string) *config.Configuration {
	cfg := &config.Configuration{}
	part := ""

	conId, ojName := false, false
	for dir != "" {
		dir, part = path.Split(strings.TrimRight(dir, "/"))
		if _, err := strconv.ParseInt(part, 10, 32); err == nil {
			cfg.CurrentContestId = part
			conId = true
		} else if strings.ToLower(part) == "cf" || strings.ToLower(part) == "codeforces" {
			cfg.OJ = "codeforces"
			ojName = true
		}
		if conId && ojName {
			cfg.Workspace = dir
			return cfg
		}
	}
	if !conId || !ojName || cfg.Workspace == "" {
		return nil
	}
	return cfg
}
