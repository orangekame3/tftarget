/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
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

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Terraform destroy, interactively select resource to destroy with target option",
	Long:  "Terraform destroy, interactively select resource to destroy with target option",
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := exec.Command("terraform", "plan", "-destroy", "-no-color").CombinedOutput()
		if err != nil {
			color.Red.Println(string(out))
			return err
		}
		resources := ExtractResourceNames(out)
		selectedResources := make([]string, 0)
		prompt := &survey.MultiSelect{
			Message: "Select resources to target destroy:",
			Options: resources,
		}
		if err := survey.AskOne(prompt, &selectedResources, survey.WithPageSize(25)); err != nil {
			if err == terminal.InterruptErr {
				log.Fatal("interrupted")
			}
		}
		targets := SliceToString(DropAction(selectedResources))

		buf := TargetCommand("destroy", targets)
		destroyCmd := exec.Command("sh", "-c", buf.String())
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Color("green")
		s.Start()
		destroyCmd.Stdout = os.Stdout
		destroyCmd.Run()
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
	rootCmd.AddCommand(destroyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// destroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
