package main

import (
	"fmt"
	"strings"
)

func versionPart(parts []string, i int) int {
	var result int
	if i >= len(parts) {
		return 0
	}
	value := parts[i]
	if d := strings.IndexByte(value, '-'); d >= 0 {
		value = value[:d]
	}
	if _, err := fmt.Sscanf(value, "%d", &result); err != nil {
		return 0
	}
	return result
}

func compareVersions(a, b string) int {
	a = strings.TrimPrefix(strings.TrimSpace(a), "v")
	b = strings.TrimPrefix(strings.TrimSpace(b), "v")
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	length := len(aParts)
	if len(bParts) > length {
		length = len(bParts)
	}
	for i := 0; i < length; i++ {
		ai := versionPart(aParts, i)
		bi := versionPart(bParts, i)
		if ai < bi {
			return -1
		}
		if ai > bi {
			return 1
		}
	}
	return 0
}
