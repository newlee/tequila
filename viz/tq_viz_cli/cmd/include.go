package cmd

import (
	"github.com/spf13/cobra"
	. "github.com/newlee/tequila/viz"
)

var includeCmd = &cobra.Command{
	Use:   "include",
	Short: "include dependencies of source code",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		result := ParseCodeDir(cmd.Flag("source").Value.String())
		result.ToDot(cmd.Flag("output").Value.String())
	},
}

func init() {
	rootCmd.AddCommand(includeCmd)

	 includeCmd.Flags().StringP("source", "s", "", "source code directory")
	 includeCmd.Flags().StringP("output", "o", "dep.dot", "output dot file name")
}
