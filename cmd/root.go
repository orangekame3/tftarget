/*
Copyright © 2023 orangekame3 <miya.org.0309@gmai.com>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tftarget",
	Short: "interactivity select resource to plan/appply with target option",
	Long: `tftarget is a CLI library for Terraform plan/apply with target option.
You can interactivity select resource to plan/appply with target option.
`,
}

var S = spinner.New(spinner.CharSets[14], 100*time.Millisecond)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}
