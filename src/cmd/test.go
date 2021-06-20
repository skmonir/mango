package cmd

import (
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/skmonir/mango/src/config"
	"github.com/skmonir/mango/src/system"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Tests a problem solution",
	Run: func(cmd *cobra.Command, args []string) {
		argc := len(args)
		if argc < 1 {
			ansi.Println(color.RedString("too few argument for 'test' command"))
		} else if argc > 1 {
			ansi.Println(color.RedString("too much argument for 'test' command"))
		} else {
			cfg, err := config.GetConfig()
			if err != nil {
				ansi.Println(color.RedString(err.Error()))
				return
			}
			if cfg.Workspace == "" {
				ansi.Println(color.RedString("workspace directory path missing. run 'mango configure' to set workspace"))
				ansi.Println(color.WhiteString("guessing workspace, platform and contestID from current directory"))
			}
			if err := system.RunTest(cfg, args[0]); err != nil {
				ansi.Println(color.RedString(err.Error()))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
