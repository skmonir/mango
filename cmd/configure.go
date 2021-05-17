package cmd

import (
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/skmonir/mango/config"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Sets current contest",
	Run: func(cmd *cobra.Command, args []string) {
		argc := len(args)
		if argc > 0 {
			ansi.Println(color.RedString("no argument required for 'configure' command"))
		} else {
			if err := config.Configure(); err != nil {
				ansi.Println(color.RedString(err.Error()))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
