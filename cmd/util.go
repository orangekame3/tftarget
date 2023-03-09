package cmd

import (
	"bytes"
	"regexp"
	"strings"
)

func ExtractResourceNames(input []byte) []string {
	re := regexp.MustCompile(`#\s([^(\n]*)(\n|$)`)

	matches := re.FindAllSubmatch(input, -1)

	var results []string
	for _, match := range matches {
		results = append(results, string(match[1]))
	}

	return results
}
func DropAction(strs []string) []string {
	var result []string
	for _, s := range strs {
		s = strings.TrimSpace(s)
		parts := strings.Split(s, " ")
		if len(parts) > 0 {
			result = append(result, parts[0])
		}
	}
	return result
}

func SliceToString(slice []string) string {
	var buffer bytes.Buffer
	if len(slice) == 1 {
		buffer.WriteString(slice[0])
		return buffer.String()
	}
	buffer.WriteString(`{`)
	for i, item := range slice {
		buffer.WriteString(item)
		if i < len(slice)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString(`}`)
	return buffer.String()
}
