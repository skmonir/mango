package intel

import (
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/skmonir/mango/src/config"
)

func GuessWorkspace(cfg *config.Configuration) bool {
	ansi.Println(color.RedString("workspace directory path missing. run 'mango configure' to set workspace"))

	guessed := WhereAmI()
	if guessed == nil {
		ansi.Println(color.RedString("couldn't guess workspace"))
		return false
	}
	ansi.Println(color.WhiteString("guessing workspace"))
	cfg.Workspace = guessed.Workspace
	cfg.CurrentContestId = guessed.CurrentContestId
	cfg.OJ = guessed.OJ
	ansi.Println(color.GreenString("guessed! Workspace:%s, OJ:%s, Contest:%s", cfg.Workspace, cfg.OJ, cfg.CurrentContestId))
	return true
}

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
