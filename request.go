package main

import (
	"net/http"
	"strings"
)

// provide a map of TYPE:key => value which is extracted from the provided request
func extractDataFromRequest(mapping map[string]string, r *http.Request) map[string]string {
	var data = make(map[string]string)
	for key, value := range mapping {
		s := strings.Split(key, ":")
		rType, keyName := strings.ToLower(s[0]), s[1]
		if rType == "get" {
			val := r.URL.Query()[keyName]
			if len(val) > 0 {
				data[value] = val[0]
			}
		} else if rType == "post" {
			panic("post variables not implemented yet")
		}
	}
	return data
}
