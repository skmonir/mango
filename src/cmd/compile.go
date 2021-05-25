package cmd

import (
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/skmonir/mango/src/config"
	"github.com/skmonir/mango/src/system"
	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compiles a source",
	Run: func(cmd *cobra.Command, args []string) {
		argc := len(args)
		if argc < 1 {
			ansi.Println(color.RedString("too few argument for 'compile' command"))
		} else if argc > 1 {
			ansi.Println(color.RedString("too much argument for 'compile' command"))
		} else {
			cfg, err := config.GetConfig()
			if err != nil {
				ansi.Println(color.RedString(err.Error()))
				return
			}
			if cfg.Workspace == "" {
				ansi.Println(color.RedString("workspace directory path missing. run 'mango configure' to set workspace"))
				return
			}
			if err := system.Compile(cfg, args[0]); err != nil {
				ansi.Println(color.RedString(err.Error()))
			} else {
				ansi.Println(color.GreenString("compiled successfully"))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
