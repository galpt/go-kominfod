package main

import (
	"errors"
	"strings"
)

func strSplit(s, sep string) ([]string, error) {
	result := strings.Split(s, sep)

	if len(result) == 1 {
		return nil, errors.New("delimiter not found")
	}

	return result, nil
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func chkIdx(elems []string, v string) int {
	for idx, s := range elems {
		if v == s {
			return idx
		}
	}
	return 0
}
