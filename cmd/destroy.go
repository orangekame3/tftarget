/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
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

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Terraform destroy, interactively select resource to destroy with target option",
	Long:  "Terraform destroy, interactively select resource to destroy with target option",
	RunE: func(cmd *cobra.Command, args []string) error {

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " loading ..."
		s.Color("green")
		s.Start()
		options, err := ExecutePlan()
		if err != nil {
			return fmt.Errorf("plan :%w", err)
		}
		s.Stop()

		selectedResources := make([]string, 0, 100)
		prompt := &survey.MultiSelect{
			Message: "Select resources to target destroy:",
			Options: options,
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
		buf := TargetCommand("destroy", targets)
		destroyCmd := exec.Command("sh", "-c", buf.String())
		s.Restart()
		destroyCmd.Stdout = os.Stdout
		destroyCmd.Run()
		s.Stop()
		if IsYes(bufio.NewReader(os.Stdin)) {
			return Confirm(buf).Run()
		}
		color.Green.Println("destroy did not executed")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// destroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
