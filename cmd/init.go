package cmd

import (
	"github.com/pietdaniel/ddmon/lib"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialze the directory for ddmon usage",
	Long: `Constructs the appropriate files and folders needed for ddmon usage.

example:
  ddmon init      # initialzes the current directory
  ddmon init $DIR # initializes the given directory
`,
	Run: func(cmd *cobra.Command, args []string) {
		lib.Run(args)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
