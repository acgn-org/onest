package tools

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func ConvertPatternRegexpString(s, regStr, pattern string) (string, error) {
	reg, err := regexp.Compile(regStr)
	if err != nil {
		return "", fmt.Errorf("compiling regular expression failed: %w", err)
	}
	return ConvertPatternRegexp(s, reg, pattern), nil
}

func ConvertPatternRegexp(s string, reg *regexp.Regexp, pattern string) string {
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
	})
}
