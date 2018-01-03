package cmd

import (
	//"fmt"
	"fmt"
	. "github.com/newlee/tequila/viz"
	"github.com/spf13/cobra"
	"strings"
)

var collCmd *cobra.Command = &cobra.Command{
	Use:   "coll",
	Short: "full collaboration grpahh",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		source := cmd.Flag("source").Value.String()
		filter := cmd.Flag("filter").Value.String()
		ignore := cmd.Flag("ignore").Value.String()
		result := ParseColl(source, filter)

		if cmd.Flag("mergePackage").Value.String() == "true" {
			result = result.MergeHeaderFile(MergePackageFunc)
		}
		if cmd.Flag("entryPoints").Value.String() == "true" {
			entryPoints := result.EntryPoints(MergePackageFunc)
			for _, cf := range entryPoints {
				fmt.Println(cf)
			}
			return
		}

		if cmd.Flag("fanInFanOut").Value.String() == "true" {
			fans := result.SortedByFan(MergePackageFunc)
			fmt.Println("Name\tTotal\tFan-In\tFan-Out")
			for _, fan := range fans {
				fmt.Printf("%s\t%v\t%v\t%v\n", fan.Name, fan.FanIn+fan.FanOut, fan.FanIn, fan.FanOut)
			}
			return
		}

		ignores := strings.Split(ignore, ",")
		var nodeFilter = func(key string) bool {
			for _, f := range ignores {
				if key == f {
					return true
				}
			}
			return false
		}
		result.ToDot(cmd.Flag("output").Value.String(), ".", nodeFilter)
	},
}

func init() {
	rootCmd.AddCommand(collCmd)

	collCmd.Flags().StringP("source", "s", "", "source code directory")
	collCmd.Flags().StringP("filter", "f", "coll__graph.dot", "dot file filter")
	collCmd.Flags().StringP("ignore", "i", "main.cpp,main", "ignore")
	collCmd.Flags().StringP("output", "o", "dep.dot", "output dot file name")
	collCmd.Flags().BoolP("entryPoints", "E", false, "list entry points")
	collCmd.Flags().BoolP("fanInFanOut", "F", false, "sorted fan-in and fan-out")
	collCmd.Flags().BoolP("mergePackage", "P", false, "merge package/folder for include dependencies")
}
