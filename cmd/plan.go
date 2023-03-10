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
		S.Suffix = " loading ..."
		S.Color("green")
		S.Start()
		options, err := ExecutePlan("")
		if err != nil && !IsNotFound(err) {
			return fmt.Errorf("plan :%w", err)
		}
		if IsNotFound(err) {
			return nil
		}
		S.Stop()

		selected := make([]string, 0, 100)
		if err := survey.AskOne(&survey.MultiSelect{Message: "Select resources to target destroy:", Options: options}, &selected, survey.WithPageSize(25)); err != nil {
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
		S.Restart()
		if err := TargetCmd(GenTargetCmd("plan", SliceToString(DropAction(selected)))).Run(); err != nil {
			return err
		}
		S.Stop()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
