package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mango",
	Short: "CP task parser and tester",
	Long:  `A command line interface(CLI) application for competitive programming task parsing and testing`,
	Run: func(cmd *cobra.Command, args []string) {
		ansi.Println(color.CyanString("Welcome to 'mango' - task parser and tester. Run 'mango --help' to see all available commands"))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
