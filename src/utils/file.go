package utils

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"time"
)

// ReadFile read file into lines
func ReadFile(filename string) (line []string, err error) {
	line = make([]string, 0)
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = append(line, scanner.Text())
	}

	return
}

// WriteFile write struct slice to file
func WriteFile(filename string, lines []string) (err error) {
	dir := path.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm|os.ModeAppend)
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, line := range lines {
		if _, err := file.WriteString(fmt.Sprintln(line)); err != nil {
			return err
		}
	}

	return
}

// ResolveOutputDir resolve output dir
func ResolveOutputDir() string {
	pwd, _ := os.Getwd()
	now := time.Now().Format("20060102")
	return path.Join(pwd, "output", now)
}
