/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Create a Terraform plan",
	Long:  `Create a Terraform plan and display the result.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := exec.Command("terraform", "plan", "-no-color").CombinedOutput()
		if err != nil {
			color.Red.Println(string(out))
			return err
		}
		resources := ExtractResourceNames(out)
		selectedResources := make([]string, 0)
		prompt := &survey.MultiSelect{
			Message: "Select resources to apply:",
			Options: resources,
		}
		survey.AskOne(prompt, &selectedResources, survey.WithPageSize(25))
		targets := SliceToString(DropAction(selectedResources))

		var buffer bytes.Buffer
		buffer.WriteString("terraform")
		buffer.WriteString(" plan")
		buffer.WriteString(" -target=")
		buffer.WriteString(targets)
		planCmd := exec.Command("sh", "-c", buffer.String())
		fmt.Println("Running command:", planCmd.String())
		planCmd.Stdout = os.Stdout
		planCmd.Stderr = os.Stderr
		return planCmd.Run()
	},
}

func init() {
	rootCmd.AddCommand(planCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// planCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// planCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
