package cmd

import (
	"fmt"
	. "github.com/newlee/tequila/viz"
	"github.com/spf13/cobra"
	"strings"
)

var includeCmd = &cobra.Command{
	Use:   "include",
	Short: "include dependencies of source code",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		result := ParseCodeDir(cmd.Flag("source").Value.String())
		var mergeFunc = func(input string) string {
			tmp := strings.Split(input, ".")
			if len(tmp) > 1 {
				return strings.Join(tmp[0:len(tmp)-1], ".")
			}
			return input
		}

		if cmd.Flag("findCrossRefs").Value.String() == "true" {
			crossRefs := result.FindCrossRef(mergeFunc)
			for _, cf := range crossRefs {
				fmt.Println(cf)
			}
			return
		}

		if cmd.Flag("mergeHeader").Value.String() == "true" {
			result = result.MergeHeaderFile(mergeFunc)
		}

		result.ToDot(cmd.Flag("output").Value.String())
	},
}

func init() {
	rootCmd.AddCommand(includeCmd)

	includeCmd.Flags().StringP("source", "s", "", "source code directory")
	includeCmd.Flags().StringP("output", "o", "dep.dot", "output dot file name")
	includeCmd.Flags().BoolP("findCrossRefs", "F", false, "find cross references")
	includeCmd.Flags().BoolP("mergeHeader", "M", false, "merge header file to same cpp file")
}
