/*
Copyright © 2023 orangekame3 <miya.org.0309@gmai.com>
*/
package cmd

import (
	"fmt"

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
		s.Start()
		options, err := executePlan(cmd, "")
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
		if err := targetCmd(genTargetCmd(cmd, "plan", slice2String(dropAction(selected)))).Run(); err != nil {
			return err
		}
		s.Stop()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.Flags().IntP("parallel", "p", 10, "Limit the number of concurrent operations as Terraform walks the graph. Defaults to 10.")
	planCmd.Flags().IntP("items", "i", 25, "Check box item size")
}
