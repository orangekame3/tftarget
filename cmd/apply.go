/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Create a Terraform apply",
	Long:  `Create a Terraform apply and display the result.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := exec.Command("terraform", "plan", "-no-color").CombinedOutput()
		if err != nil {
			color.Red.Println(string(out))
			return err
		}
		resources := ExtractResourceNames(out)
		selectedResources := make([]string, 0)
		prompt := &survey.MultiSelect{
			Message: "Select resources to target apply:",
			Options: resources,
		}
		survey.AskOne(prompt, &selectedResources, survey.WithPageSize(25))
		targets := SliceToString(DropAction(selectedResources))

		var buffer bytes.Buffer
		buffer.WriteString("terraform")
		buffer.WriteString(" apply")
		buffer.WriteString(" -target=")
		buffer.WriteString(targets)
		applyCmd := exec.Command("sh", "-c", buffer.String())
		applyCmd.Stdout = os.Stdout
		applyCmd.Stderr = os.Stderr
		applyCmd.Run()

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Do you want to perform these actions? ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text != "yes" {
			return nil
		}

		buffer.WriteString(" -auto-approve")
		confirm := exec.Command("sh", "-c", buffer.String())
		confirm.Stdout = os.Stdout
		confirm.Stderr = os.Stderr
		return confirm.Run()
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
