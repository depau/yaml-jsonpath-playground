package main

import (
	"strings"
	"testing"
)

func TestEvaluation(t *testing.T) {
	testCases := []struct {
		name     string
		yamlDoc  string
		jsonPath string
		expected string
	}{
		{
			name:     "Extract names from a list",
			yamlDoc:  "- name: foo\n- name: bar",
			jsonPath: "[*].name",
			expected: "- foo\n- bar",
		},
		{
			name: "Extract values from multiple ConfigMaps",
			yamlDoc: `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: foo
  namespace: foo
data:
  value: foo
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: bar
  namespace: bar
data:
  value: bar`,
			jsonPath: "data.value",
			expected: `---
# ConfigMap: foo (namespace: foo)
foo
---
# ConfigMap: bar (namespace: bar)
bar`,
		},
		{
			name: "Extract value from a single ConfigMap",
			yamlDoc: `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: foo
  namespace: foo
data:
  value: foo`,
			jsonPath: "data.value",
			expected: "foo",
		},
		{
			name:     "Invalid YAML",
			yamlDoc:  "a: b: c",
			jsonPath: "$.a",
			expected: "YAML Error: yaml: mapping values are not allowed in this context",
		},
		{
			name:     "Invalid JSONPath",
			yamlDoc:  "a: b",
			jsonPath: "$..",
			expected: "JSONPath Error: child name or array access or filter missing after recursive descent at position 3, following \"$..\"",
		},
		{
			name:     "No results",
			yamlDoc:  "a: b",
			jsonPath: "$.c",
			expected: "",
		},
		{
			name:     "Empty YAML",
			yamlDoc:  "",
			jsonPath: "$.a",
			expected: "",
		},
		{
			name:     "Multi-document with no results on one",
			yamlDoc:  `---`,
			jsonPath: "$.a",
			expected: ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := playgroundEvaluate(tc.yamlDoc, tc.jsonPath)
			if strings.TrimSpace(actual) != strings.TrimSpace(tc.expected) {
				t.Errorf(`For jsonpath '%s', expected:
"""
%s
"""
got:
"""
%s
"""`, tc.jsonPath, tc.expected, actual)
			}
		})
	}
}
