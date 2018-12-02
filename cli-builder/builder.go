package clibuilder

import (
	"regexp"
	"strings"
)

// Parse parse the cli pattern
func Parse(pattern string, data map[string]string) string {
	r, err := regexp.Compile(`\{\s*(.*?)\s*\}`)
	if err != nil {
		panic(err)
	}
	replaced := r.ReplaceAllStringFunc(pattern, func(s string) string {
		part := r.FindStringSubmatch(s)[1]
		return part

	})
	return strings.TrimSpace(replaced)
}
