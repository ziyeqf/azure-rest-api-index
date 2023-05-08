package main

import (
	"bufio"
	"bytes"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type TagInfo struct {
	InputFile []string `yaml:"input-file"`
}

func SpecListFromReadmeMD(b []byte) ([]string, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	specSet := map[string]struct{}{}
	var isEnter bool
	var ymlContent string
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "```yaml $(tag) ==") {
			isEnter = true
			continue
		}
		if trimmedLine == "```" {
			var info TagInfo
			if err := yaml.Unmarshal([]byte(ymlContent), &info); err != nil {
				return nil, fmt.Errorf("decoding yaml %q: %v", ymlContent, err)
			}
			for _, p := range info.InputFile {
				specSet[p] = struct{}{}
			}
			// rest the states
			isEnter = false
			ymlContent = ""
			continue
		}
		if !isEnter {
			continue
		}
		ymlContent += line + "\n"
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %v", err)
	}
	var out []string
	for p := range specSet {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i] < out[j]
	})
	return out, nil
}
