package cmd

import (
	"strconv"

	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/skmonir/mango/config"
	"github.com/spf13/cobra"
)

var setContestCmd = &cobra.Command{
	Use:   "setc",
	Short: "Sets current contest",
	Run: func(cmd *cobra.Command, args []string) {
		argc := len(args)
		if argc < 1 {
			ansi.Println(color.RedString("too few argument for 'setc' command"))
		} else if argc > 1 {
			ansi.Println(color.RedString("too much argument for 'setc' command"))
		} else {
			if _, err := strconv.Atoi(args[0]); err != nil {
				ansi.Println(color.RedString("contest id not valid"))
			} else if err := config.SetContest(args[0]); err != nil {
				ansi.Println(color.RedString(err.Error()))
			} else {
				ansi.Println(color.New(color.FgGreen).Sprintf("current contest is set to %v", args[0]))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(setContestCmd)
}
