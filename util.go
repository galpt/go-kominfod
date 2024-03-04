package main

import (
	"errors"
	"strings"
)

// custom strings.Split
func strSplit(s, sep string) ([]string, error) {
	result := strings.Split(s, sep)

	if len(result) == 1 {
		return nil, errors.New("delimiter not found")
	}

	return result, nil
}
