package clibuilder

import (
	"fmt"
	"regexp"
	"strings"
)

// Parse parse the cli pattern
func Parse(pattern string, data map[string]string) string {
	r, _ := regexp.Compile(`\{\s*(--){0,2}(\w*)\s*(=?)\s*(\w*)\s*\}`)
	replaced := r.ReplaceAllStringFunc(pattern, func(s string) string {
		parts := r.FindStringSubmatch(s)
		flag, key, equals, defaultVal := parts[1], parts[2], parts[3], parts[4]

		var isFlag = false
		if len(flag) > 0 {
			isFlag = true
		}

		if val, ok := data[key]; ok {
			if isFlag {
				if len(equals) == 0 {
					return flag + key
				}
				return fmt.Sprintf("%s%s=%s", flag, key, val)
			} else {
				return val
			}
		} else if len(defaultVal) > 0 {
			// is there a default
			if isFlag {
				return fmt.Sprintf("%s%s=%s", flag, key, defaultVal)
			} else {
				return defaultVal
			}
		}

		// remove pattern
		return ""

	})
	return strings.TrimSpace(replaced)
}
