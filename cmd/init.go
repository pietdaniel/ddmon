package cmd

import (
	"github.com/pietdaniel/ddmon/lib"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
