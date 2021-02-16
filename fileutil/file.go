package fileutil

import (
	"bufio"
	"os"
)

func ReadLines(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var res []string
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	return res, nil
}

func MustReadLines(filename string) []string {
	paths, err := ReadLines(filename)
	if err != nil {
		panic(err)
	}
	return paths
}