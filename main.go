package main

import (
	"bytes"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("evaluate", js.FuncOf(evaluate))
	<-c
}

func evaluate(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return "Invalid number of arguments"
	}
	yamlDoc := args[0].String()
	jsonPath := args[1].String()

	var n yaml.Node
	if err := yaml.Unmarshal([]byte(yamlDoc), &n); err != nil {
		return fmt.Sprintf("YAML Error: %s", err)
	}

	path, err := yamlpath.NewPath(jsonPath)
	if err != nil {
		return fmt.Sprintf("JSONPath Error: %s", err)
	}

	results, err := path.Find(&n)
	if err != nil {
		return fmt.Sprintf("Evaluation Error: %s", err)
	}

	out := []string{}
	for _, a := range results {
		b, err := encode(a)
		if err != nil {
			return fmt.Sprintf("Output Encoding Error: %s", err)
		}
		out = append(out, b)
	}

	return strings.Join(out, "---\n")
}

func encode(a *yaml.Node) (string, error) {
	var buf bytes.Buffer
	e := yaml.NewEncoder(&buf)
	defer e.Close()
	e.SetIndent(2)

	if err := e.Encode(a); err != nil {
		return "", err
	}

	return buf.String(), nil
}
