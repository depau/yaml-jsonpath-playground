package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

func playgroundEvaluate(yamlDoc string, jsonPath string) string {
	path, err := yamlpath.NewPath(jsonPath)
	if err != nil {
		return fmt.Sprintf("JSONPath Error: %s", err)
	}

	decoder := yaml.NewDecoder(strings.NewReader(yamlDoc))

	var documents []yaml.Node
	var outputs []string
	for {
		var n yaml.Node
		if err := decoder.Decode(&n); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Sprintf("YAML Error: %s", err)
		}
		documents = append(documents, n)
	}

	for _, n := range documents {
		if len(n.Content) == 0 || (len(n.Content) == 1 && n.Content[0].Tag == "!!null") {
			outputs = append(outputs, "")
			continue
		}

		var parts []string

		if len(documents) > 1 {
			if comment := genComment(n); comment != nil {
				parts = append(parts, *comment)
			}
		}

		results, err := path.Find(&n)
		if err != nil {
			parts = append(parts, fmt.Sprintf("Evaluation Error: %s", err))
		} else if len(results) > 0 {
			var nodeToEncode *yaml.Node
			if len(results) == 1 && !strings.HasPrefix(jsonPath, "[*]") {
				nodeToEncode = results[0]
			} else {
				nodeToEncode = &yaml.Node{
					Kind:    yaml.SequenceNode,
					Content: results,
				}
			}
			b, err := encode(nodeToEncode)
			if err != nil {
				return fmt.Sprintf("Output Encoding Error: %s", err)
			}
			parts = append(parts, b)
		}

		outputs = append(outputs, strings.Join(parts, "\n"))
	}

	if len(documents) == 0 {
		return ""
	} else if len(documents) == 1 {
		return outputs[0]
	}

	hasContent := false
	for _, out := range outputs {
		if out != "" {
			hasContent = true
			break
		}
	}

	if !hasContent {
		return ""
	}

	result := strings.Join(outputs, "\n---\n")
	result = "---\n" + result

	return result
}

func genComment(n yaml.Node) *string {
	if kind, ok := findScalar(&n, "kind"); ok {
		if name, ok := findScalar(&n, "metadata", "name"); ok {
			comment := fmt.Sprintf("# %s: %s", kind, name)
			if namespace, ok := findScalar(&n, "metadata", "namespace"); ok {
				comment += fmt.Sprintf(" (namespace: %s)", namespace)
			}
			return &comment
		}
	}
	return nil
}

func encode(a *yaml.Node) (string, error) {
	var buf bytes.Buffer
	e := yaml.NewEncoder(&buf)
	defer e.Close()
	e.SetIndent(2)

	if err := e.Encode(a); err != nil {
		return "", err
	}

	return strings.TrimSuffix(buf.String(), "\n"), nil
}

func findScalar(node *yaml.Node, path ...string) (string, bool) {
	if node == nil {
		return "", false
	}
	if len(path) == 0 {
		if node.Kind == yaml.ScalarNode {
			return node.Value, true
		}
		return "", false
	}

	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		return findScalar(node.Content[0], path...)
	}

	if node.Kind != yaml.MappingNode {
		return "", false
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Value == path[0] {
			return findScalar(node.Content[i+1], path[1:]...)
		}
	}
	return "", false
}
