package cmd

import (
	//"fmt"

	. "github.com/newlee/tequila/viz"
	"github.com/spf13/cobra"
	//"strings"
	"strings"
)

var callCmd *cobra.Command = &cobra.Command{
	Use:   "call",
	Short: "icall grpah",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		source := cmd.Flag("source").Value.String()
		filter := cmd.Flag("filter").Value.String()
		result := ParseICallGraph(source, filter)
		clusterResults := []string{
			"com.lenovo.awakens.model.product.Configuration.",
		}
		var nodeFilter = func(key string) bool {
			//if strings.HasPrefix(key, "com.lenovo.awakens.model.product.Configuration.") {
			for _, cs := range clusterResults {
				if strings.Contains(key, cs) {
					return false
				}
			}
			//}

			return true
		}
		result.ToDot(cmd.Flag("output").Value.String(), ".", nodeFilter)
		result.ToDataSet("", ".", nodeFilter)
	},
}

func init() {
	rootCmd.AddCommand(callCmd)

	callCmd.Flags().StringP("source", "s", "", "source code directory")
	callCmd.Flags().StringP("filter", "f", "coll__graph.dot", "dot file filter")
	callCmd.Flags().StringP("output", "o", "dep.dot", "output dot file name")
}
