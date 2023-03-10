/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
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
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " loading ..."
		s.Color("green")
		s.Start()
		out, err := exec.Command("terraform", "plan", "-no-color").CombinedOutput()
		if err != nil {
			color.Red.Println(string(out))
			return fmt.Errorf("plan :%w", err)
		}
		resources := make([]string, 0, 100)
		resources = append(resources, color.Red.Sprintf("%s", "exit (cancel terraform plan)"))
		resources = append(resources, ExtractResourceNames(out)...)
		selectedResources := make([]string, 0)
		s.Stop()
		prompt := &survey.MultiSelect{
			Message: "Select resources to target plan:",
			Options: resources,
		}
		if err := survey.AskOne(prompt, &selectedResources, survey.WithPageSize(25)); err != nil {
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
		targets := SliceToString(DropAction(selectedResources))
		buf := TargetCommand("plan", targets)
		planCmd := exec.Command("sh", "-c", buf.String())
		s.Restart()
		planCmd.Stdout = os.Stdout
		planCmd.Stderr = os.Stderr
		if err := planCmd.Run(); err != nil {
			return fmt.Errorf("target plan :%w", err)
		}
		s.Stop()
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
