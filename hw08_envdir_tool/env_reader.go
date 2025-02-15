package main

import (
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	content, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envList := make(Environment)
	for _, file := range content {
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), "=") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		info, err := file.Info()
		if err != nil {
			return nil, err
		}

		dataString := string(data)
		newLinePos := strings.Index(dataString, "\n")

		if newLinePos >= 0 {
			dataString = dataString[:newLinePos]
		}

		envList[file.Name()] = EnvValue{
			Value:      strings.ReplaceAll(strings.TrimRight(dataString, " \t"), "\x00", "\n"),
			NeedRemove: info.Size() == 0,
		}
	}

	return envList, nil
}
