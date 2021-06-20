package intel

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/skmonir/mango/src/config"
)

func WhereAmI() (*config.Configuration, error) {
	dir, err := os.Getwd()
	return whereAmI(dir), err
}

func whereAmI(dir string) *config.Configuration {
	cfg := &config.Configuration{}
	part := ""

	for dir != "" {
		fmt.Println(dir)
		dir, part = path.Split(strings.TrimRight(dir, "/"))
		if _, err := strconv.ParseInt(part, 10, 32); err == nil {
			cfg.CurrentContestId = part
		} else if strings.ToLower(part) == "cf" || strings.ToLower(part) == "codeforces" {
			cfg.OJ = "codeforces"
		}
	}
	return cfg
}
