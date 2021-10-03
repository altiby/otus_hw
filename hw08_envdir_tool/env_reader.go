package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func readDirFiles(dir string) (map[string][]byte, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed read dir %s, %w", dir, err)
	}

	retMap := make(map[string][]byte, len(files))

	for _, file := range files {
		if !file.IsDir() {
			fileData, err := readFileContent(dir + "/" + file.Name())
			if err != nil {
				return nil, fmt.Errorf("failed read file %s, %w", file.Name(), err)
			}
			retMap[file.Name()] = fileData
		}
	}
	return retMap, nil
}

func readFileContent(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed read file %s, %w", fileName, err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	data, _, err := reader.ReadLine()
	if errors.Is(err, io.EOF) {
		return []byte{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed read file data %s, %w", fileName, err)
	}

	return data, nil
}

func normalizeData(data []byte) string {
	retStr := bytes.TrimRight(
		bytes.TrimRight(
			bytes.ReplaceAll(data, []byte{0}, []byte{'\n'}), " "), "\t")
	return string(retStr)
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here
	dirData, err := readDirFiles(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment, len(dirData))
	for k, v := range dirData {
		nd := normalizeData(v)
		if len(nd) == 0 {
			env[k] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
		} else {
			env[k] = EnvValue{
				Value:      nd,
				NeedRemove: false,
			}
		}
	}
	return env, nil
}
