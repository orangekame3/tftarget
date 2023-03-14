/*
Copyright © 2023 orangekame3 <miya.org.0309@gmai.com>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
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
		s.Start()
		options, err := executePlan(cmd, "-destroy")
		if err != nil && !IsNotFound(err) {
			return fmt.Errorf("plan :%w", err)
		}
		if IsNotFound(err) {
			return nil
		}
		s.Stop()

		selected := make([]string, 0, 100)
		items, _ := cmd.Flags().GetInt("items")
		if err := survey.AskOne(&survey.MultiSelect{Message: "Select resources to target destroy:", Options: options}, &selected, survey.WithPageSize(items)); err != nil {
			return fmt.Errorf("select resource :%w", err)
		}
		if len(selected) == 0 {
			color.Green.Println("resource not seleced")
			return nil
		}
		if slices.Contains(selected, color.Red.Sprintf("%s", "exit (cancel terraform plan)")) {
			color.Green.Println("exit seleced")
			return nil
		}
		s.Restart()
		buf := genTargetCmd(cmd, "destroy", slice2String(dropAction(selected)))
		targetCmd(buf).Run()
		s.Stop()
		if isYes(bufio.NewReader(os.Stdin)) {
			return confirm(buf).Run()
		}
		color.Green.Println("destroy did not executed")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
	destroyCmd.Flags().IntP("parallel", "p", 10, "Limit the number of concurrent operations as Terraform walks the graph. Defaults to 10.")
	destroyCmd.Flags().IntP("items", "i", 25, "Check box item size")
}
