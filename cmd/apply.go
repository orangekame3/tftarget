/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Terraform apply, interactively select resource to apply with target option",
	Long:  "Terraform apply, interactively select resource to apply with target option",
	RunE: func(cmd *cobra.Command, args []string) error {
		options, err := ExecutePlan()
		if err != nil {
			return fmt.Errorf("plan :%w", err)
		}
		S.Stop()

		selectedResources := make([]string, 0, 100)
		if err := survey.AskOne(&survey.MultiSelect{Message: "Select resources to target destroy:", Options: options}, &selectedResources, survey.WithPageSize(25)); err != nil {
			return fmt.Errorf("select resource :%w", err)
		}
		if len(selectedResources) == 0 {
			color.Green.Println("resource not seleced")
			return nil
		}
		if slices.Contains(selectedResources, color.Red.Sprintf("%s", "exit (cancel terraform plan)")) {
			color.Green.Println("exit seleced")
			return nil
		}
		buf := TargetCommand("apply", SliceToString(DropAction(selectedResources)))
		applyCmd := exec.Command("sh", "-c", buf.String())
		S.Restart()
		applyCmd.Stdout = os.Stdout
		applyCmd.Run()
		S.Stop()
		if IsYes(bufio.NewReader(os.Stdin)) {
			return Confirm(buf).Run()
		}
		color.Green.Println("apply did not executed")
		return nil
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
