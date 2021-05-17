package cmd

import (
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version number of 'mango'",
	Run: func(cmd *cobra.Command, args []string) {
		ansi.Println(color.CyanString("'mango' task parser and tester -- v0.9"))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
