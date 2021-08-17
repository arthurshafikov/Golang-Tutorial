package main

import (
	"io/ioutil"
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
	envMap := make(Environment)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		fileName := f.Name()
		if strings.ToUpper(fileName) != fileName || strings.Contains(fileName, "=") {
			// ignore not ENV typed files
			continue
		}

		file, err := ioutil.ReadFile(dir + "/" + fileName)
		if err != nil {
			return nil, err
		}

		var value string
		needRemove := false
		if f.Size() < 1 {
			needRemove = true
		} else {
			value = strings.Split(string(file), "\n")[0]
			value = strings.ReplaceAll(value, "\x00", "\n")
			value = strings.TrimRight(value, " \t")
		}

		envValue := EnvValue{
			Value:      value,
			NeedRemove: needRemove,
		}

		envMap[fileName] = envValue
	}
	return envMap, nil
}
