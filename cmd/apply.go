/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
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
			Message: "Select resources to target apply:",
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
		buf := TargetCommand("apply", targets)
		applyCmd := exec.Command("sh", "-c", buf.String())
		s.Restart()
		applyCmd.Stdout = os.Stdout
		applyCmd.Run()
		s.Stop()
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Do you want to perform these actions? ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text != "yes" {
			return nil
		}
		buf.WriteString(" -auto-approve")
		confirm := exec.Command("sh", "-c", buf.String())
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
