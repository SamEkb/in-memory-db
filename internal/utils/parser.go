package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseTime(duration time.Duration) time.Time {
	res := time.Now().Add(duration)
	return res
}

func ParseSize(sizeStr string) (int, error) {
	sizeStr = strings.ToUpper(sizeStr)

	var multiplier int
	switch {
	case strings.HasSuffix(sizeStr, "KB"):
		multiplier = 1024
		sizeStr = strings.TrimSuffix(sizeStr, "KB")
	case strings.HasSuffix(sizeStr, "MB"):
		multiplier = 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "MB")
	case strings.HasSuffix(sizeStr, "GB"):
		multiplier = 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "GB")
	case strings.HasSuffix(sizeStr, "B"):
		multiplier = 1
		sizeStr = strings.TrimSuffix(sizeStr, "B")
	default:
		return 0, fmt.Errorf("unknown size suffix in %s", sizeStr)
	}

	value, err := strconv.Atoi(sizeStr)
	if err != nil {
		return 0, fmt.Errorf("invalid size value: %v", err)
	}

	return value * multiplier, nil
}
