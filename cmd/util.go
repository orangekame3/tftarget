package cmd

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gookit/color"
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

func TargetCommand(cmd, target string) bytes.Buffer {
	var buf bytes.Buffer
	buf.WriteString("terraform ")
	buf.WriteString(cmd)
	buf.WriteString(" -target=")
	buf.WriteString(target)
	return buf
}

func IsYes(reader *bufio.Reader) bool {
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text) == "yes"
}

func Confirm(buf bytes.Buffer) *exec.Cmd {
	buf.WriteString(" -auto-approve")
	confirm := exec.Command("sh", "-c", buf.String())
	confirm.Stdout = os.Stdout
	confirm.Stderr = os.Stderr
	return confirm
}

func ExecutePlan() ([]string, error) {
	out, err := exec.Command("terraform", "plan", "-no-color").CombinedOutput()
	if err != nil {
		color.Red.Println(string(out))
		return nil, err
	}
	resources := ExtractResourceNames(out)
	if len(resources) == 0 {
		color.Green.Println(string(out))
		return nil, nil
	}
	options := make([]string, 0, 100)
	options = append(options, color.Red.Sprintf("%s", "exit (cancel terraform plan)"))
	return append(options, resources...), nil
}
