package clibuilder

import (
	"testing"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		pattern  string
		args     map[string]string
		expected string
	}{
		// no flags
		{"path", map[string]string{}, "path"},
		// binary
		{"path {--flag}", map[string]string{}, "path"},
		{"path {--flag}", map[string]string{"flag": "someval"}, "path --flag"},
		{"path { --flag  }", map[string]string{"flag": "someval"}, "path --flag"},
		{"path {--flag=}", map[string]string{"flag": "someval"}, "path --flag=someval"},
		{"path {flag}", map[string]string{"flag": "someval"}, "path someval"},
		{"path {flag} | pipe {to}", map[string]string{"flag": "someval"}, "path someval | pipe"},
		{"path {flag} | pipe {to}", map[string]string{"flag": "someval", "to": "another"}, "path someval | pipe another"},
		{"path {flag} | pipe {to=hi}", map[string]string{"flag": "someval"}, "path someval | pipe hi"},
	}

	for _, test := range tests {
		if output := Parse(test.pattern, test.args); output != test.expected {
			t.Error("Test failed: {} inputted, {} expected, received {}", test.pattern, test.expected, output)
		}
	}
}
