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
	return res, scanner.Err()
}

func MustReadLines(filename string) []string {
	paths, err := ReadLines(filename)
	if err != nil {
		panic(err)
	}
	return paths
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
