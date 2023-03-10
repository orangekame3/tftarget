/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Terraform plan, interactively select resource to plan with target option",
	Long:  "Terraform plan, interactively select resource to plan with target option",
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
		buf := TargetCommand("plan", SliceToString(DropAction(selectedResources)))
		planCmd := exec.Command("sh", "-c", buf.String())
		S.Restart()
		planCmd.Stdout = os.Stdout
		if err := planCmd.Run(); err != nil {
			return fmt.Errorf("target plan :%w", err)
		}
		S.Stop()
		return nil
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
