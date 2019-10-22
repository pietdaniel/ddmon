package cmd

import (
	"github.com/pietdaniel/ddmon/lib"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates all templates into terraform files",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		lib.Generate(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	var TargetDir string
	var SourceDir string
	var TemplateDir string

	generateCmd.Flags().StringVarP(&TargetDir, "target-dir", "t", "./output", "Target directory to write to")
	generateCmd.Flags().StringVarP(&SourceDir, "source-dir", "s", "./", "Target directory to read from")
	generateCmd.Flags().StringVarP(&TemplateDir, "template-dir", "p", "./templates", "Target directory to read templates from")

	generateCmd.MarkFlagRequired("target-dir")
	generateCmd.MarkFlagRequired("source-dir")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
