/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/briandowns/spinner"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Terraform apply, interactively select resource to apply with target option",
	Long:  "Terraform apply, interactively select resource to apply with target option",
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
		if err := survey.AskOne(prompt, &selectedResources, survey.WithPageSize(25)); err != nil {
			if err == terminal.InterruptErr {
				log.Fatal("interrupted")
			}
		}
		targets := SliceToString(DropAction(selectedResources))

		var buffer bytes.Buffer
		buffer.WriteString("terraform")
		buffer.WriteString(" apply")
		buffer.WriteString(" -target=")
		buffer.WriteString(targets)
		applyCmd := exec.Command("sh", "-c", buffer.String())
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
		s.Color("green")
		s.Start()
		//applyCmd.Stderr = os.Stderr
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
