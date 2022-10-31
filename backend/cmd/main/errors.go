package main

import (
	"strings"
)

func ehSkip(err error, msg ...string) {
	if err == nil {
		return
	}
	if len(msg) > 0 {
		LOG.Println(err, msg)
	} else {
		LOG.Println(err)
	}
}

func ehIsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "not found")
}
