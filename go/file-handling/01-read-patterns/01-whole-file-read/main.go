package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const MAX_FILE_SIZE = 1024 * 1024

var ErrTooBig = errors.New("File to big")

func ReadConfig(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ReadConfig: %w", err)
	}

	defer file.Close()

	buf, err := io.ReadAll(io.LimitReader(file, MAX_FILE_SIZE+1))
	if err != nil {
		return nil, fmt.Errorf("ReadConfig: %w", err)
	}

	if int64(len(buf)) > MAX_FILE_SIZE {
		return nil, fmt.Errorf("ReadConfig: %w", ErrTooBig)
	}

	return buf, nil
}

func main() {
	filePath := "README.md"
	bytes, err := ReadConfig(filePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(string(bytes))
}
