/*
Copyright © 2023 orangekame3 <miya.org.0309@gmai.com>
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

func extractResource(input []byte) []string {
	re := regexp.MustCompile(`#\s([^(\n]*)(\n|$)`)

	matches := re.FindAllSubmatch(input, -1)

	var results []string
	for _, match := range matches {
		results = append(results, string(match[1]))
	}

	return results
}
func dropAction(strs []string) []string {
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

func slice2String(slice []string) string {
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

func genTargetCmd(cmd *cobra.Command, action, target string) bytes.Buffer {
	var buf bytes.Buffer
	buf.WriteString("terraform ")
	buf.WriteString(action)
	buf.WriteString(" -target=")
	buf.WriteString(target)
	p, _ := cmd.Flags().GetInt("parallel")
	buf.WriteString(fmt.Sprintf(" --parallelism=%d", p))
	return buf
}

func isYes(reader *bufio.Reader) bool {
	color.Red.Print("Enter a value: ")
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text) == "yes"
}

func confirm(buf bytes.Buffer) *exec.Cmd {
	buf.WriteString(" -auto-approve")
	confirm := exec.Command("sh", "-c", buf.String())
	confirm.Stdout = os.Stdout
	confirm.Stderr = os.Stderr
	return confirm
}

func executePlan(cmd *cobra.Command, option string) ([]string, error) {
	p, _ := cmd.Flags().GetInt("parallel")
	planCmd := exec.Command("terraform", "plan", "-no-color", fmt.Sprintf("--parallelism=%d", p))
	if option != "" {
		planCmd = exec.Command("terraform", "plan", option, "-no-color", fmt.Sprintf("--parallelism=%d", p))
	}

	out, err := planCmd.CombinedOutput()
	if err != nil {
		color.Red.Println(string(out))
		return nil, err
	}
	resources := extractResource(out)
	if len(resources) == 0 {
		color.Green.Println(string(out))
		return nil, ErrNotFound
	}
	options := make([]string, 0, 100)
	options = append(options, color.Red.Sprintf("%s", "exit (cancel terraform plan)"))
	return append(options, resources...), nil
}

func targetCmd(buf bytes.Buffer) *exec.Cmd {
	cmd := exec.Command("sh", "-c", buf.String())
	cmd.Stdout = os.Stdout
	return cmd
}
