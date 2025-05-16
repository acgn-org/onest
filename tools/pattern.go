package tools

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func ConvertWithPattern(s, regStr, pattern string) (string, error) {
	reg, err := regexp.Compile(regStr)
	if err != nil {
		return "", fmt.Errorf("error compiling regular expression: %w", err)
	}
	matches := reg.FindStringSubmatch(s)
	return os.Expand(pattern, func(s string) string {
		i, err := strconv.Atoi(s)
		if err != nil {
			return ""
		}
		if i >= 0 && i < len(matches) {
			return matches[i]
		}
		return ""
	}), nil
}
