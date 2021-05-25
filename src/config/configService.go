package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
)

func GetConfigFolderPath() string {
	// return filepath.Join(os.Getenv("HOME"), ".mango") // "MacOS"
	// return filepath.Join(os.Getenv("APPDATA"), ".mango") // "windows"
	// if os.Getenv("XDG_CONFIG_HOME") != "" {
	// 	globalSettingFolder = os.Getenv("XDG_CONFIG_HOME")
	// } else {
	// 	globalSettingFolder = filepath.Join(os.Getenv("HOME"), ".config")
	// }

	cfgPath := ""
	switch runtime.GOOS {
	case "linux":
		if os.Getenv("XDG_CONFIG_HOME") != "" {
			cfgPath = os.Getenv("XDG_CONFIG_HOME")
		} else {
			cfgPath = filepath.Join(os.Getenv("HOME"), ".mango")
		}
	case "windows":
		cfgPath = filepath.Join(os.Getenv("APPDATA"), "mango")
	case "darwin":
		cfgPath = filepath.Join(os.Getenv("HOME"), ".mango")
	default:
		cfgPath = ""
	}

	return cfgPath
}

func GetConfigFilePath() string {
	return filepath.Join(GetConfigFolderPath(), "config.json")
}

func IsConfigDirExist() bool {
	_, err := os.Stat(GetConfigFolderPath())
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateConfigDir() error {
	if !IsConfigDirExist() {
		if err := os.MkdirAll(GetConfigFolderPath(), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func GetHost(OJ string) string {
	HostMap := map[string]string{
		"CF": "https://codeforces.com",
	}

	host, ok := HostMap[OJ]
	if !ok {
		return ""
	}
	return host
}

func GetFullOnlineJudgeName(OJ string) string {
	OnlineJudgeNameMap := map[string]string{
		"CF": "codeforces",
	}

	oj, ok := OnlineJudgeNameMap[OJ]
	if !ok {
		return ""
	}
	return oj
}

func GetConfig() (Configuration, error) {
	CreateDefaultConfig()

	cfgFilePath := GetConfigFilePath()

	data, err := ioutil.ReadFile(cfgFilePath)
	if err != nil {
		return Configuration{}, err
	}

	var cfg Configuration

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Configuration{}, err
	}
	return cfg, nil
}

func SaveConfig(cfg Configuration) error {
	data, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}

	mu := sync.Mutex{}

	cfgPath := GetConfigFilePath()
	mu.Lock()
	err = ioutil.WriteFile(cfgPath, data, 0644)
	mu.Unlock()
	if err != nil {
		return err
	}

	return nil
}

func CreateDefaultConfig() error {
	if IsConfigExist() {
		return nil
	}

	CreateConfigDir()

	cfg := Configuration{
		CompilationCommand: "g++",
		CompilationArgs:    "-std=c++17",
		OJ:                 "codeforces",
		Host:               "https://codeforces.com",
	}

	if err := SaveConfig(cfg); err != nil {
		fmt.Println(err.Error())
		return errors.New("error while creating default configuration")
	}

	return nil
}

func SetOnlineJudge(OJ string) error {
	cfg, err := GetConfig()
	if err != nil {
		return err
	}

	cfg.Host = GetHost(OJ)
	cfg.OJ = GetFullOnlineJudgeName(OJ)

	if err := SaveConfig(cfg); err != nil {
		return err
	}

	return nil
}

func SetContest(contestId string) error {
	if _, err := strconv.Atoi(contestId); err != nil {
		return errors.New("contest id not valid")
	}

	cfg, err := GetConfig()
	if err != nil {
		return err
	}

	cfg.CurrentContestId = contestId

	if err := SaveConfig(cfg); err != nil {
		return err
	}

	return nil
}

func Configure() error {
	var err error

	if err = CreateDefaultConfig(); err != nil {
		return err
	}

	cfgPath := GetConfigFilePath()

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", cfgPath).Run()
	case "windows":
		exec.Command("cmd", fmt.Sprintf("/C start %v", cfgPath)).Run()
	case "darwin":
		err = exec.Command("open", cfgPath).Run()
	default:
		ansi.Println(color.New(color.FgRed).Sprintf("unsupported os"))
	}

	return err
}

func IsConfigExist() bool {
	cfgFilePath := GetConfigFilePath()
	info, err := os.Stat(cfgFilePath)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
