package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

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
